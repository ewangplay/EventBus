package redis

import (
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

type Client struct {
	ctx  *Context
	pool *redis.Pool
	sync *redsync.Redsync
}

func NewClient(ctx *Context) (*Client, error) {
	rc := &Client{}

	rc.ctx = ctx
	rc.pool = &redis.Pool{
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

	rc.sync = redsync.New([]redsync.Pool{rc.pool})

	return rc, nil
}

func (rc *Client) Close() error {
	return rc.pool.Close()
}

func (rc *Client) GetRedisSync() *redsync.Redsync {
	return rc.sync
}

func (rc *Client) Incr(key string) (uint64, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.Uint64(conn.Do("INCR", key))
}

func (rc *Client) Exists(key string) (bool, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

func (rc *Client) Get(key string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (rc *Client) Set(key string, value string) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (rc *Client) HDel(key, field string) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", key, field)
	if err != nil {
		return err
	}

	return nil
}

func (rc *Client) HExists(key, field string) (bool, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("HEXISTS", key, field))
}

func (rc *Client) HGet(key, field string) ([]byte, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("HGET", key, field))
}

func (rc *Client) HGetAll(key string) (map[string]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}

func (rc *Client) HKeys(key string) ([]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("HKEYS", key))
}

func (rc *Client) HLen(key string) (uint64, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return redis.Uint64(conn.Do("HLEN", key))
}

func (rc *Client) HSet(key, field string, value []byte) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return err
	}
	return nil
}

func (rc *Client) HSetNX(key, field string, value []byte) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSETNX", key, field, value)
	if err != nil {
		return err
	}
	return nil
}
