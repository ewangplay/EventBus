package rabbitmq

import (
	"time"

	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

var (
	gRmqCtx *Context
)

type Context struct {
	opts   *config.EBOptions
	logger i.Logger
}

func GetContext(opts *config.EBOptions, logger i.Logger) (*Context, error) {
	if gRmqCtx == nil {
		gRmqCtx = &Context{}
		gRmqCtx.opts = opts
		gRmqCtx.logger = logger
	}

	return gRmqCtx, nil
}

func (c *Context) CreateProducer() (i.Producer, error) {
	c.logger.Debug("Creating Producer ...")

	publisher, err := NewProducer(c)
	if err != nil {
		return nil, err
	}

	return publisher, err
}

func (c *Context) CreateConsumer(subject string) (i.Consumer, error) {
	c.logger.Debug("Creating Consumer ...")

	subscriber, err := NewConsumer(c, subject)
	if err != nil {
		return nil, err
	}

	return subscriber, err
}

func (c *Context) createSessionPool(ctx context.Context) chan chan AMQPSession {
	sessionPool := make(chan chan AMQPSession)

	go func() {
		sess := make(chan AMQPSession)

		for {
			select {
			case sessionPool <- sess:
			case <-ctx.Done():
				c.logger.Info("Shutting down session factory")
				return
			}

			c.logger.Debug("Dialing amqp server %s ...", c.opts.RabbitmqAddress)
			var conn *amqp.Connection
			var ch *amqp.Channel
			var err error
			for {
				conn, err = amqp.Dial(c.opts.RabbitmqAddress)
				if err != nil {
					c.logger.Error("Cannot (re)dial: %+v: %q, retry after 500 millisecond", err, c.opts.RabbitmqAddress)
					time.Sleep(500 * time.Millisecond)
					continue
				}

				ch, err = conn.Channel()
				if err != nil {
					c.logger.Error("Cannot create channel: %v, retry after 500 millisecond", err)
					conn.Close()
					time.Sleep(500 * time.Millisecond)
					continue
				}

				c.logger.Debug("Declaring the exchange name: %s, kind: %s.", c.opts.RabbitmqExchangeName, c.opts.RabbitmqExchangeKind)
				err = ch.ExchangeDeclare(c.opts.RabbitmqExchangeName, c.opts.RabbitmqExchangeKind, c.opts.RabbitmqDurable, c.opts.RabbitmqAutoDelete, false, false, nil)
				if err != nil {
					c.logger.Error("Cannot declare %s exchange named as %s: %+v, retry after 500 millisecond", c.opts.RabbitmqExchangeKind, c.opts.RabbitmqExchangeName, err)
					conn.Close()
					time.Sleep(500 * time.Millisecond)
					continue
				}

				//loop only once
				break
			}

			select {
			case sess <- AMQPSession{conn, ch}:
			case <-ctx.Done():
				c.logger.Info("Shutting down new session")
				return
			}
		}
	}()

	return sessionPool
}
