package rabbitmq

import (
	"fmt"
	"time"

	comm "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type Consumer struct {
	i.Logger
	ctx         *Context
	netCtx      context.Context
	netCtxDone  context.CancelFunc
	sessionPool chan chan AMQPSession
	messages    chan i.Message
	useSession  AMQPSession
}

func NewConsumer(ctx *Context, subject string) (*Consumer, error) {
	c := &Consumer{}
	c.Logger = ctx.logger
	c.ctx = ctx
	c.messages = make(chan i.Message)
	c.netCtx, c.netCtxDone = context.WithCancel(context.Background())
	c.sessionPool = c.ctx.createSessionPool(c.netCtx)
	go func() {
		c.subscribe(subject, c.sessionPool, c.messages)
	}()

	return c, nil
}

func (c *Consumer) Close() error {
	c.netCtxDone()
	c.useSession.Close()
	close(c.messages)
	close(c.sessionPool)
	c.Info("rabbitmq consumer stopped.")
	return nil
}

func (c *Consumer) GetMessage() <-chan i.Message {
	return c.messages
}

func (c *Consumer) subscribe(subject string, sessionPool chan chan AMQPSession, messages chan<- i.Message) {

	var err error
	var queue amqp.Queue
	var msgs <-chan amqp.Delivery

	for sessions := range sessionPool {
		session := <-sessions
		c.useSession = session

		queue, err = session.QueueDeclare(
			subject, // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			c.Error("Failed to declare a queue: %v", err)
			time.Sleep(500 * time.Millisecond)
			session.Close()
			continue
		}

		//routingKey := fmt.Sprintf("%s.*", subject)
		routingKey := fmt.Sprintf("%s.#", subject)
		err = session.QueueBind(
			queue.Name,                      // queue name
			routingKey,                      // routing key
			c.ctx.opts.RabbitmqExchangeName, // exchange
			false,
			nil)
		if err != nil {
			c.Error("Failed to bind a queue: %v", err)
			time.Sleep(500 * time.Millisecond)
			session.Close()
			continue
		}

		c.Info("Binding queue %s to exchange %s with routing key %s", queue.Name, c.ctx.opts.RabbitmqExchangeName, routingKey)

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
			c.Error("Failed to register a consumer: %v", err)
			time.Sleep(500 * time.Millisecond)
			session.Close()
			continue
		}

		c.Info("Waiting for deliver message...")

		for d := range msgs {
			c.Debug("Received message[%s: %s]", subject, d.Body)
			c.messages <- &comm.EBMessage{Subject: subject, Data: d.Body}
		}
	}
}
