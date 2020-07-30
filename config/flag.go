package config

import (
	"flag"

	"github.com/ewangplay/eventbus/utils"
)

func NewFlagSet(opts *EB_Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("eventbus", flag.ExitOnError)

	// basic options
	flagSet.Bool("version", false, "print version string")
	flagSet.String("config", "service.yaml", "path to config file")

	// eventbus flags
	flagSet.String("http-address", opts.EBHTTPAddress, "<addr>:<port> to listen on for http request")
	flagSet.String("url-pattern", opts.EBUrlPattern, "http request url regexp pattern")
	flagSet.Int("retry-policy", opts.EBRetryPolicy, "retry policy for failed event: 1 - count policy; 2 - timeout policy")
	flagSet.Int("max-retry-count", opts.EBMaxRetryCount, "max retry count for failed event")
	flagSet.Duration("retry-interval", opts.EBRetryInterval, "duration time of retry interval for failed event")
	flagSet.Duration("max-retry-timeout", opts.EBMaxRetryTimeout, "max retry timeout for failed event")

	// drivers flags
	drivers := utils.StringArray{}
	flagSet.Var(&drivers, "driver", "driver to enable(may be given multiple times)")

	// log flags
	flagSet.String("log-mode", "normal", "set log mode: normal, dev")
	flagSet.String("log-level", "info", "set log verbosity: 1: debug, 2: info, 3: warn, 4: error, 5: fatal")
	flagSet.String("log-prefix", "[EventBus] ", "log message prefix")
	flagSet.String("log-path", opts.LogPath, "path to log files")
	flagSet.Int64("log-max-size", opts.LogMaxSize, "max size per log file before rolling (megabytes)")
	flagSet.Bool("log-rotate-daily", opts.LogRotateDaily, "whether the log file rotate daily")
	flagSet.Int("log-max-age", opts.LogMaxAge, "max age per log file (days)")

	// pprof flags
	flagSet.Bool("pprof-enable", false, "whether to enable pprof debug")
	flagSet.String("pprof-address", opts.PProfAddress, "pprof debug service address")

	// redis flags
	flagSet.Bool("redis-enable", true, "whether to enable redis service")
	flagSet.Bool("redis-cluster", false, "whether to enable redis cluster")
	flagSet.String("redis-address", opts.RedisAddress, "<addr>:<port> to connect redis service")
	flagSet.String("redis-credential", opts.RedisCredential, "redis connection credential")
	flagSet.Int("redis-max-idle", opts.RedisMaxIdle, "maximum number of idle connections in the pool.")
	flagSet.Duration("redis-idle-timeout", opts.RedisIdleTimeout, "close connections after remaining idle for this duration")

	// rabbitmq flags
	flagSet.Bool("rabbitmq-enable", opts.RabbitmqEnable, "whether to enable rabbitmq service")
	flagSet.String("rabbitmq-address", opts.RabbitmqAddress, "<addr>:<port> to connect rabbitmq service")
	flagSet.String("rabbitmq-exchange-name", opts.RabbitmqExchangeName, "rabbitmq exchange name setting")
	flagSet.String("rabbitmq-exchange-kind", opts.RabbitmqExchangeKind, "rabbitmq exchange kind setting")
	flagSet.Bool("rabbitmq-durable", opts.RabbitmqDurable, "rabbitmq durable flag")
	flagSet.Bool("rabbitmq-auto-delete", opts.RabbitmqAutoDelete, "rabbitmq auto-delete flag")

	// nsq flags
	flagSet.Bool("nsq-enable", opts.NSQEnable, "whether to enable nsqd service")
	flagSet.Bool("nsq-cluster", opts.NSQCluster, "whether to enable nsqd cluster")
	flagSet.Int("nsq-max-in-flight", opts.NSQMaxInFlight, "Maximum number of messages to allow in flight")
	flagSet.String("nsq-tcp-address", opts.NSQTCPAddress, "<addr>:<port> to listen on for TCP clients")
	lookupdTCPAddrs := utils.StringArray{}
	flagSet.Var(&lookupdTCPAddrs, "nsq-lookupd-tcp-address", "lookupd TCP address (may be given multiple times)")

	return flagSet
}

