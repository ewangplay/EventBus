package solo

import (
	"sync"

	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

// Constants Define
const (
	QueueMaxSize = 10000
)

var (
	gSoloCtx *Context
)

// Context represents Solo Mode Messager Context
type Context struct {
	i.Logger
	opts *config.EBOptions
	sync.RWMutex
	Queues map[string]chan []byte
}

// GetContext returns the Solo Mode Messager Context in singleton mode
func GetContext(opts *config.EBOptions, logger i.Logger) (*Context, error) {
	if gSoloCtx == nil {
		gSoloCtx = &Context{}
		gSoloCtx.opts = opts
		gSoloCtx.Logger = logger
		gSoloCtx.Queues = make(map[string]chan []byte)
	}
	return gSoloCtx, nil
}

// CreateProducer creates Solo  Mode Messager Producer instance
func (c *Context) CreateProducer() (i.Producer, error) {
	c.Debug("Creating Producer ...")

	publisher, err := NewProducer(c)
	if err != nil {
		return nil, err
	}

	return publisher, err
}

// CreateConsumer creates Solo Mode Messager Consumer instance
func (c *Context) CreateConsumer(topic string) (i.Consumer, error) {
	c.Debug("Creating Consumer ...")

	subscriber, err := NewConsumer(c, topic)
	if err != nil {
		return nil, err
	}

	return subscriber, err
}
