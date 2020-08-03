package adapter

import (
	"fmt"
	"strconv"

	"github.com/ewangplay/eventbus/adapter/redis"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

// Counter struct define, implment Counter interface
type Counter struct {
	i.Logger
	redisCtx    *redis.Context
	redisClient *redis.Client
}

// NewCounter ...
func NewCounter(opts *config.EBOptions, logger i.Logger) (*Counter, error) {
	idcounter := &Counter{logger, nil, nil}

	redisCtx, err := redis.GetContext(opts, logger)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisCtx.CreateRedisClient()
	if err != nil {
		return nil, err
	}

	idcounter.redisCtx = redisCtx
	idcounter.redisClient = redisClient

	return idcounter, nil
}

// Close ...
func (c *Counter) Close() error {
	c.Info("Id counter will close")
	return c.redisClient.Close()
}

// NewEventID ...
func (c *Counter) NewEventID() (string, error) {
	if c.redisClient == nil {
		return "", fmt.Errorf("redis client instance invalid")
	}

	key := fmt.Sprintf("EVENTBUS:EVENT:ID")
	id, err := c.redisClient.Incr(key)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(id, 10), nil
}
