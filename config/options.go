package config

import (
	"time"

	"github.com/ewangplay/eventbus/i"
)

type EB_Options struct {
	ServiceName string

	// eventbus options
	EBHTTPAddress     string        `flag:"http-address" cfg:"eventbus_http_address"`
	EBUrlPattern      string        `flag:"url-pattern" cfg:"eventbus_url_pattern"`
	EBRetryPolicy     int           `flag:"retry-policy" cfg:"eventbus_retry_policy"`
	EBMaxRetryCount   int           `flag:"max-retry-count" cfg:"eventbus_max_retry_count"`
	EBRetryInterval   time.Duration `flag:"retry-interval" cfg:"eventbus_retry_interval"`
	EBMaxRetryTimeout time.Duration `flag:"max-retry-timeout" cfg:"eventbus_max_retry_timeout"`

	// driver options
	Drivers []string `flag:"driver"`

	// log options
	Logger         i.ILogger
	LogMode        string   `flag:"log-mode"`
	LogLevel       LogLevel `flag:"log-level"`
	LogPrefix      string   `flag:"log-prefix"`
	LogPath        string   `flag:"log-path"`
	LogMaxSize     int64    `flag:"log-max-size"`
	LogRotateDaily bool     `flag:"log-rotate-daily"`
	LogMaxAge      int      `flag:"log-max-age"`

	// pprof options
	PProfEnable  bool   `flag:"pprof-enable"`
	PProfAddress string `flag:"pprof-address"`

	// redis options
	RedisEnable      bool          `flag:"redis-enable"`
	RedisCluster     bool          `flag:"redis-cluster"`
	RedisAddress     string        `flag:"redis-address"`
	RedisCredential  string        `flag:"redis-credential"`
	RedisMaxIdle     int           `flag:"redis-max-idle"`
	RedisIdleTimeout time.Duration `flag:"redis-idle-timeout"`

	// rabbitmq options
	RabbitmqEnable       bool   `flag:"rabbitmq-enable"`
	RabbitmqExchangeName string `flag:"rabbitmq-exchange-name"`
	RabbitmqExchangeKind string `flag:"rabbitmq-exchange-kind"`
	RabbitmqAddress      string `flag:"rabbitmq-address"`
	RabbitmqDurable      bool   `flag:"rabbitmq-durable"`
	RabbitmqAutoDelete   bool   `flag:"rabbitmq-auto-delete"`

	// nsq options
	NSQEnable  bool `flag:"nsq-enable"`
	NSQCluster bool `flag:"nsq-cluster"`
	NSQMaxInFlight int `flag:"nsq-max-in-flight"`
	NSQTCPAddress          string   `flag:"nsq-tcp-address"`
	NSQLookupdTCPAddresses []string `flag:"nsq-lookupd-tcp-address" cfg:"nsq_lookupd_tcp_addresses"`
}

func NewOptions() *EB_Options {

	return &EB_Options{
		ServiceName: "eventbus",

		EBHTTPAddress:     "0.0.0.0:8091",
		EBUrlPattern:      "/(v\\d+)/(\\w+)",
		EBRetryPolicy:     1,
		EBMaxRetryCount:   10,
		EBRetryInterval:   5 * time.Second,
		EBMaxRetryTimeout: 2 * time.Hour,

		Drivers: make([]string, 0),

		LogMode:        "normal",
		LogLevel:       2,
		LogPrefix:      "[EventBus] ",
		LogPath:        "/opt/ewangplay/eventbus/log",
		LogMaxSize:     100,
		LogRotateDaily: false,
		LogMaxAge:      30,

		PProfEnable:  false,
		PProfAddress: "0.0.0.0:6060",

		RedisEnable:      true,
		RedisCluster:     false,
		RedisAddress:     "127.0.0.1:7001",
		RedisMaxIdle:     30,
		RedisIdleTimeout: 1 * time.Minute,

		RabbitmqEnable:       false,
		RabbitmqExchangeName: "eventbus",
		RabbitmqExchangeKind: "topic",
		RabbitmqAddress:      "amqp://guest:111111@127.0.0.1:5672/",
		RabbitmqDurable:      true,
		RabbitmqAutoDelete:   true,

		NSQEnable:      true,
		NSQCluster:     false,
		NSQMaxInFlight: 1000,
		NSQTCPAddress:   "127.0.0.1:4150",
		NSQLookupdTCPAddresses: make([]string, 0),
	}
}
