package nsq

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

var (
	gNsqCtx *Context
)

type Context struct {
	i.Logger
	opts *config.EBOptions
}

func GetContext(opts *config.EBOptions, logger i.Logger) (*Context, error) {
	if gNsqCtx == nil {
		gNsqCtx = &Context{}
		gNsqCtx.opts = opts
		gNsqCtx.Logger = logger
	}
	return gNsqCtx, nil
}

func (c *Context) CreateProducer() (i.Producer, error) {
	c.Debug("Creating Producer ...")

	publisher, err := NewProducer(c)
	if err != nil {
		return nil, err
	}

	return publisher, err
}

func (c *Context) CreateConsumer(topic string) (i.Consumer, error) {
	c.Debug("Creating Consumer ...")

	subscriber, err := NewConsumer(c, topic)
	if err != nil {
		return nil, err
	}

	return subscriber, err
}
