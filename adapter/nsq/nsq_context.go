package nsq

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

var (
	g_nsqCtx *NSQContext
)

type NSQContext struct {
	i.ILogger
	opts *config.EB_Options
}

func GetNSQContext(opts *config.EB_Options, logger i.ILogger) (*NSQContext, error) {
	if g_nsqCtx == nil {
		g_nsqCtx = &NSQContext{}
		g_nsqCtx.opts = opts
		g_nsqCtx.ILogger = logger
	}
	return g_nsqCtx, nil
}

func (this *NSQContext) CreateProducer() (i.IProducer, error) {
	this.Debug("Creating NSQProducer ...")

	publisher, err := NewNSQProducer(this)
	if err != nil {
		return nil, err
	}

	return publisher, err
}

func (this *NSQContext) CreateConsumer(topic string) (i.IConsumer, error) {
	this.Debug("Creating NSQConsumer ...")

	subscriber, err := NewNSQConsumer(this, topic)
	if err != nil {
		return nil, err
	}

	return subscriber, err
}
