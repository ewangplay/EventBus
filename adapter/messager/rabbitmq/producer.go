package rabbitmq

import (
	"github.com/ewangplay/eventbus/i"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type AMQPSession struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s AMQPSession) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

type Producer struct {
	i.Logger
	ctx         *Context
	netCtx      context.Context
	netCtxDone  context.CancelFunc
	sessionPool chan chan AMQPSession
	messages    chan i.Message
	useSession  AMQPSession
}

func NewProducer(ctx *Context) (*Producer, error) {
	p := &Producer{}
	p.Logger = ctx.logger
	p.ctx = ctx
	p.messages = make(chan i.Message)
	p.netCtx, p.netCtxDone = context.WithCancel(context.Background())
	p.sessionPool = p.ctx.createSessionPool(p.netCtx)
	go func() {
		p.publish(p.sessionPool, p.messages)
	}()

	return p, nil
}

func (p *Producer) Close() error {
	p.netCtxDone()
	p.useSession.Close()
	close(p.messages)
	close(p.sessionPool)
	p.Info("rabbitmq producer server stopped.")
	return nil
}

func (p *Producer) Publish(msg i.Message) error {
	p.messages <- msg
	return nil
}

// publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func (p *Producer) publish(sessionPool chan chan AMQPSession, messages <-chan i.Message) {
	var (
		running bool
		reading = messages
		pending = make(chan i.Message, 1)
	)

	for sessions := range sessionPool {
		session := <-sessions
		p.useSession = session

		confirm := make(chan amqp.Confirmation, 1)

		// producer confirms for p channel/connection
		if err := session.Confirm(false); err != nil {
			p.Debug("producer confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			session.NotifyPublish(confirm)
		}

		p.Debug("publishing...")

		var message i.Message

	CURR_SESSION:
		for {
			select {
			case ack := <-confirm:
				if !ack.Ack {
					p.Debug("nack message %d, subject: %s, body: %s", ack.DeliveryTag,
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
				p.Debug("Publishing message [%s:%s]", message.GetSubject(), message.GetData())
				err := session.Publish(p.ctx.opts.RabbitmqExchangeName,
					message.GetSubject(),
					false,
					false,
					amqp.Publishing{
						Body: []byte(message.GetData()),
					})
				// Retry failed delivery on the next session
				if err != nil {
					p.Error("Failed to publish message [%s:%s], retry on the next session.", message.GetSubject(), message.GetData())
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
