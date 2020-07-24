package rabbitmq

import (
	"time"

	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

var (
	g_rbtCtx *RabbitMQContext
)

type RabbitMQContext struct {
	opts   *config.EB_Options
	logger i.ILogger
}

func GetRabbitMQContext(opts *config.EB_Options, logger i.ILogger) (*RabbitMQContext, error) {
	if g_rbtCtx == nil {
		g_rbtCtx = &RabbitMQContext{}
		g_rbtCtx.opts = opts
		g_rbtCtx.logger = logger
	}

	return g_rbtCtx, nil
}

func (this *RabbitMQContext) CreateProducer() (i.IProducer, error) {
	this.logger.Debug("Creating RabbitMQProducer ...")

	publisher, err := NewRabbitMQProducer(this)
	if err != nil {
		return nil, err
	}

	return publisher, err
}

func (this *RabbitMQContext) CreateConsumer(subject string) (i.IConsumer, error) {
	this.logger.Debug("Creating RabbitMQConsumer ...")

	subscriber, err := NewRabbitMQConsumer(this, subject)
	if err != nil {
		return nil, err
	}

	return subscriber, err
}

func (this *RabbitMQContext) createSessionPool(ctx context.Context) chan chan AMQP_Session {
	session_pool := make(chan chan AMQP_Session)

	go func() {
		sess := make(chan AMQP_Session)

		for {
			select {
			case session_pool <- sess:
			case <-ctx.Done():
				this.logger.Info("Shutting down session factory")
				return
			}

			this.logger.Debug("Dialing amqp server %s ...", this.opts.RabbitmqAddress)
			var conn *amqp.Connection
			var ch *amqp.Channel
			var err error
			for {
				conn, err = amqp.Dial(this.opts.RabbitmqAddress)
				if err != nil {
					this.logger.Error("Cannot (re)dial: %+v: %q, retry after 500 millisecond", err, this.opts.RabbitmqAddress)
					time.Sleep(500 * time.Millisecond)
					continue
				}

				ch, err = conn.Channel()
				if err != nil {
					this.logger.Error("Cannot create channel: %v, retry after 500 millisecond", err)
					conn.Close()
					time.Sleep(500 * time.Millisecond)
					continue
				}

				this.logger.Debug("Declaring the exchange name: %s, kind: %s.", this.opts.RabbitmqExchangeName, this.opts.RabbitmqExchangeKind)
				err = ch.ExchangeDeclare(this.opts.RabbitmqExchangeName, this.opts.RabbitmqExchangeKind, this.opts.RabbitmqDurable, this.opts.RabbitmqAutoDelete, false, false, nil)
				if err != nil {
					this.logger.Error("Cannot declare %s exchange named as %s: %+v, retry after 500 millisecond", this.opts.RabbitmqExchangeKind, this.opts.RabbitmqExchangeName, err)
					conn.Close()
					time.Sleep(500 * time.Millisecond)
					continue
				}

				//loop only once
				break
			}

			select {
			case sess <- AMQP_Session{conn, ch}:
			case <-ctx.Done():
				this.logger.Info("Shutting down new session")
				return
			}
		}
	}()

	return session_pool
}
