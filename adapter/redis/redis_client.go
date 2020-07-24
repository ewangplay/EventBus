package redis

import (
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

type RedisClient struct {
	ctx  *RedisContext
	pool *redis.Pool
	sync *redsync.Redsync
}

func NewRedisClient(ctx *RedisContext) (*RedisClient, error) {
	this := &RedisClient{}

	this.ctx = ctx
	this.pool = &redis.Pool{
		MaxIdle:     ctx.opts.RedisMaxIdle,
		IdleTimeout: ctx.opts.RedisIdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ctx.opts.RedisAddress)
			if err != nil {
				return nil, err
			}
			if ctx.opts.RedisCredential != "" {
				if _, err := c.Do("AUTH", ctx.opts.RedisCredential); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	this.sync = redsync.New([]redsync.Pool{this.pool})

	return this, nil
}

func (this *RedisClient) Close() error {
	return this.pool.Close()
}

func (this *RedisClient) GetRedisSync() *redsync.Redsync {
	return this.sync
}

func (this *RedisClient) Incr(key string) (uint64, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Uint64(conn.Do("INCR", key))
}

func (this *RedisClient) Exists(key string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

func (this *RedisClient) Get(key string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (this *RedisClient) Set(key string, value string) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (this *RedisClient) HDel(key, field string) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", key, field)
	if err != nil {
		return err
	}

	return nil
}

func (this *RedisClient) HExists(key, field string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("HEXISTS", key, field))
}

func (this *RedisClient) HGet(key, field string) ([]byte, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("HGET", key, field))
}

func (this *RedisClient) HGetAll(key string) (map[string]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}

func (this *RedisClient) HKeys(key string) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("HKEYS", key))
}

func (this *RedisClient) HLen(key string) (uint64, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Uint64(conn.Do("HLEN", key))
}

func (this *RedisClient) HSet(key, field string, value []byte) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisClient) HSetNX(key, field string, value []byte) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSETNX", key, field, value)
	if err != nil {
		return err
	}
	return nil
}
