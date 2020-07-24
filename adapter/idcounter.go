package adapter

import (
	"fmt"
	"strconv"

	"github.com/ewangplay/eventbus/adapter/redis"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

type IDCounter struct {
	i.ILogger
	redisCtx    *redis.RedisContext
	redisClient *redis.RedisClient
}

func NewIdCounter(opts *config.EB_Options, logger i.ILogger) (*IDCounter, error) {
	idcounter := &IDCounter{logger, nil, nil}

	redisCtx, err := redis.GetRedisContext(opts, logger)
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

func (this *IDCounter) Close() error {
	this.Info("Id counter will close")
	return this.redisClient.Close()
}

func (this *IDCounter) NewEventId() (string, error) {
	if this.redisClient == nil {
		return "", fmt.Errorf("redis client instance invalid")
	}

	key := fmt.Sprintf("EVENTBUS:EVENT:ID")
	id, err := this.redisClient.Incr(key)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(id, 10), nil
}
