package redis

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

var (
	g_redisContext *RedisContext
)

type RedisContext struct {
	opts   *config.EB_Options
	logger i.ILogger
}

func GetRedisContext(opts *config.EB_Options, logger i.ILogger) (*RedisContext, error) {
	if g_redisContext == nil {
		g_redisContext = &RedisContext{}
		g_redisContext.opts = opts
		g_redisContext.logger = logger
	}
	return g_redisContext, nil
}

func (r *RedisContext) CreateRedisClient() (*RedisClient, error) {
	client, err := NewRedisClient(r)
	if err != nil {
		return nil, err
	}

	return client, err
}
