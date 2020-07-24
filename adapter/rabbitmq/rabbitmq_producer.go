package rabbitmq

import (
	"github.com/ewangplay/eventbus/i"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type AMQP_Session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s AMQP_Session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

type RabbitMQProducer struct {
	i.ILogger
	ctx          *RabbitMQContext
	netCtx       context.Context
	netCtxDone   context.CancelFunc
	session_pool chan chan AMQP_Session
	messages     chan i.IMessage
	use_session  AMQP_Session
}

func NewRabbitMQProducer(ctx *RabbitMQContext) (*RabbitMQProducer, error) {
	this := &RabbitMQProducer{}
	this.ILogger = ctx.logger
	this.ctx = ctx
	this.messages = make(chan i.IMessage)
	this.netCtx, this.netCtxDone = context.WithCancel(context.Background())
	this.session_pool = this.ctx.createSessionPool(this.netCtx)
	go func() {
		this.publish(this.session_pool, this.messages)
	}()

	return this, nil
}

func (this *RabbitMQProducer) Close() error {
	this.netCtxDone()
	this.use_session.Close()
	close(this.messages)
	close(this.session_pool)
	this.Info("rabbitmq producer server stopped.")
	return nil
}

func (this *RabbitMQProducer) Publish(msg i.IMessage) error {
	this.messages <- msg
	return nil
}

// publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func (this *RabbitMQProducer) publish(session_pool chan chan AMQP_Session, messages <-chan i.IMessage) {
	var (
		running bool
		reading = messages
		pending = make(chan i.IMessage, 1)
	)

	for sessions := range session_pool {
		session := <-sessions
		this.use_session = session

		confirm := make(chan amqp.Confirmation, 1)

		// producer confirms for this channel/connection
		if err := session.Confirm(false); err != nil {
			this.Debug("producer confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			session.NotifyPublish(confirm)
		}

		this.Debug("publishing...")

		var message i.IMessage

	CURR_SESSION:
		for {
			select {
			case ack := <-confirm:
				if !ack.Ack {
					this.Debug("nack message %d, subject: %s, body: %s", ack.DeliveryTag,
						message.GetSubject(), message.GetData())
					if ack.DeliveryTag > 0 {
						//requeue current message
						pending <- message
						session.Close()
						break CURR_SESSION
					}
				} else {
					reading = messages
				}

			case message = <-pending:
				this.Debug("Publishing message [%s:%s]", message.GetSubject(), message.GetData())
				err := session.Publish(this.ctx.opts.RabbitmqExchangeName,
					message.GetSubject(),
					false,
					false,
					amqp.Publishing{
						Body: []byte(message.GetData()),
					})
				// Retry failed delivery on the next session
				if err != nil {
					this.Error("Failed to publish message [%s:%s], retry on the next session.", message.GetSubject(), message.GetData())
					pending <- message
					session.Close()
					break CURR_SESSION
				}

			case message, running = <-reading:
				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- message
				reading = nil
			}
		}
	}
}
