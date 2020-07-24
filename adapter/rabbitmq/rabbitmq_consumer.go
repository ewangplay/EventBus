package rabbitmq

import (
	"fmt"
	"time"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type RabbitMQConsumer struct {
	i.ILogger
	ctx          *RabbitMQContext
	netCtx       context.Context
	netCtxDone   context.CancelFunc
	session_pool chan chan AMQP_Session
	messages     chan i.IMessage
	use_session  AMQP_Session
}

func NewRabbitMQConsumer(ctx *RabbitMQContext, subject string) (*RabbitMQConsumer, error) {
	this := &RabbitMQConsumer{}
	this.ILogger = ctx.logger
	this.ctx = ctx
	this.messages = make(chan i.IMessage)
	this.netCtx, this.netCtxDone = context.WithCancel(context.Background())
	this.session_pool = this.ctx.createSessionPool(this.netCtx)
	go func() {
		this.subscribe(subject, this.session_pool, this.messages)
	}()

	return this, nil
}

func (this *RabbitMQConsumer) Close() error {
	this.netCtxDone()
	this.use_session.Close()
	close(this.messages)
	close(this.session_pool)
	this.Info("rabbitmq consumer stopped.")
	return nil
}

func (this *RabbitMQConsumer) GetMessage() <-chan i.IMessage {
	return this.messages
}

func (this *RabbitMQConsumer) subscribe(subject string, session_pool chan chan AMQP_Session, messages chan<- i.IMessage) {

	var err error
	var queue amqp.Queue
	var msgs <-chan amqp.Delivery

	for sessions := range session_pool {
		session := <-sessions
		this.use_session = session

		queue, err = session.QueueDeclare(
			subject, // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			this.Error("Failed to declare a queue: %v", err)
			time.Sleep(500 * time.Millisecond)
			session.Close()
			continue
		}

		//routing_key := fmt.Sprintf("%s.*", subject)
		routing_key := fmt.Sprintf("%s.#", subject)
		err = session.QueueBind(
			queue.Name,                         // queue name
			routing_key,                        // routing key
			this.ctx.opts.RabbitmqExchangeName, // exchange
			false,
			nil)
		if err != nil {
			this.Error("Failed to bind a queue: %v", err)
			time.Sleep(500 * time.Millisecond)
			session.Close()
			continue
		}

		this.Info("Binding queue %s to exchange %s with routing key %s", queue.Name, this.ctx.opts.RabbitmqExchangeName, routing_key)

		msgs, err = session.Consume(
			queue.Name, // queue
			"",         // consumer
			true,       // auto ack
			false,      // exclusive
			false,      // no local
			false,      // no wait
			nil,        // args
		)
		if err != nil {
			this.Error("Failed to register a consumer: %v", err)
			time.Sleep(500 * time.Millisecond)
			session.Close()
			continue
		}

		this.Info("Waiting for deliver message...")

		for d := range msgs {
			this.Debug("Received message[%s: %s]", subject, d.Body)
			this.messages <- &c.EB_Message{Subject: subject, Data: d.Body}
		}
	}
}
