package redis

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

var (
	gRedisCtx *Context
)

type Context struct {
	opts   *config.EBOptions
	logger i.Logger
}

func GetContext(opts *config.EBOptions, logger i.Logger) (*Context, error) {
	if gRedisCtx == nil {
		gRedisCtx = &Context{}
		gRedisCtx.opts = opts
		gRedisCtx.logger = logger
	}
	return gRedisCtx, nil
}

func (r *Context) CreateRedisClient() (*Client, error) {
	client, err := NewClient(r)
	if err != nil {
		return nil, err
	}

	return client, err
}
