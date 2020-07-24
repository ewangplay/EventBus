package config

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/utils"
)

type tlsRequiredOption int

func (t *tlsRequiredOption) Set(s string) error {
	s = strings.ToLower(s)
	if s == "tcp-https" {
		*t = c.TLSRequiredExceptHTTP
		return nil
	}
	required, err := strconv.ParseBool(s)
	if required {
		*t = c.TLSRequired
	} else {
		*t = c.TLSNotRequired
	}
	return err
}

func (t *tlsRequiredOption) Get() interface{} { return int(*t) }

func (t *tlsRequiredOption) String() string {
	return strconv.FormatInt(int64(*t), 10)
}

func (t *tlsRequiredOption) IsBoolFlag() bool { return true }

type tlsMinVersionOption uint16

func (t *tlsMinVersionOption) Set(s string) error {
	s = strings.ToLower(s)
	switch s {
	case "":
		return nil
	case "ssl3.0":
		*t = tls.VersionSSL30
	case "tls1.0":
		*t = tls.VersionTLS10
	case "tls1.1":
		*t = tls.VersionTLS11
	case "tls1.2":
		*t = tls.VersionTLS12
	default:
		return fmt.Errorf("unknown tlsVersionOption %q", s)
	}
	return nil
}

func (t *tlsMinVersionOption) Get() interface{} { return uint16(*t) }

func (t *tlsMinVersionOption) String() string {
	return strconv.FormatInt(int64(*t), 10)
}

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
	flagSet.String("log-level", "info", "set log verbosity: debug, info, warn, error, or fatal")
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
	flagSet.Bool("rabbitmq-enable", false, "whether to enable rabbitmq service")
	flagSet.String("rabbitmq-address", opts.RabbitmqAddress, "<addr>:<port> to connect rabbitmq service")
	flagSet.String("rabbitmq-exchange-name", opts.RabbitmqExchangeName, "rabbitmq exchange name setting")
	flagSet.String("rabbitmq-exchange-kind", opts.RabbitmqExchangeKind, "rabbitmq exchange kind setting")
	flagSet.Bool("rabbitmq-durable", opts.RabbitmqDurable, "rabbitmq durable flag")
	flagSet.Bool("rabbitmq-auto-delete", opts.RabbitmqAutoDelete, "rabbitmq auto-delete flag")

	// nsq flags
	flagSet.Bool("nsq-enable", true, "whether to enable nsqd service")
	flagSet.Bool("nsq-cluster", false, "whether to enable nsqd cluster")

	flagSet.Int("nsq-max-in-flight", opts.NSQMaxInFlight, "Maximum number of messages to allow in flight")

	flagSet.String("nsq-https-address", opts.NSQHTTPSAddress, "<addr>:<port> to listen on for HTTPS clients")
	flagSet.String("nsq-http-address", opts.NSQHTTPAddress, "<addr>:<port> to listen on for HTTP clients")
	flagSet.String("nsq-tcp-address", opts.NSQTCPAddress, "<addr>:<port> to listen on for TCP clients")
	flagSet.String("nsq-broadcast-address", opts.NSQBroadcastAddress, "address that will be registered with lookupd (defaults to the OS hostname)")
	authHTTPAddresses := utils.StringArray{}
	flagSet.Var(&authHTTPAddresses, "nsq-auth-http-address", "<addr>:<port> to query auth server (may be given multiple times)")
	lookupdTCPAddrs := utils.StringArray{}
	flagSet.Var(&lookupdTCPAddrs, "nsq-lookupd-tcp-address", "lookupd TCP address (may be given multiple times)")
	flagSet.Duration("http-client-connect-timeout", opts.HTTPClientConnectTimeout, "timeout for HTTP connect")
	flagSet.Duration("http-client-request-timeout", opts.HTTPClientRequestTimeout, "timeout for HTTP request")

	// diskqueue options
	flagSet.String("data-path", opts.DataPath, "path to store disk-backed messages")
	flagSet.Int64("mem-queue-size", opts.MemQueueSize, "number of messages to keep in memory (per topic/channel)")
	flagSet.Int64("max-bytes-per-file", opts.MaxBytesPerFile, "number of bytes per diskqueue file before rolling")
	flagSet.Int64("sync-every", opts.SyncEvery, "number of messages per diskqueue fsync")
	flagSet.Duration("sync-timeout", opts.SyncTimeout, "duration of time per diskqueue fsync")

	// msg and command options
	flagSet.Duration("msg-timeout", opts.MsgTimeout, "default duration to wait before auto-requeing a message")
	flagSet.Duration("max-msg-timeout", opts.MaxMsgTimeout, "maximum duration before a message will timeout")
	flagSet.Int64("max-msg-size", opts.MaxMsgSize, "maximum size of a single message in bytes")
	flagSet.Duration("max-req-timeout", opts.MaxReqTimeout, "maximum requeuing timeout for a message")
	flagSet.Int64("max-body-size", opts.MaxBodySize, "maximum size of a single command body")

	// client overridable configuration options
	flagSet.Duration("max-heartbeat-interval", opts.MaxHeartbeatInterval, "maximum client configurable duration of time between client heartbeats")
	flagSet.Int64("max-rdy-count", opts.MaxRdyCount, "maximum RDY count for a client")
	flagSet.Int64("max-output-buffer-size", opts.MaxOutputBufferSize, "maximum client configurable size (in bytes) for a client output buffer")
	flagSet.Duration("max-output-buffer-timeout", opts.MaxOutputBufferTimeout, "maximum client configurable duration of time between flushing to a client")

	// statsd integration options
	flagSet.String("statsd-address", opts.StatsdAddress, "UDP <addr>:<port> of a statsd daemon for pushing stats")
	flagSet.Duration("statsd-interval", opts.StatsdInterval, "duration between pushing to statsd")
	flagSet.Bool("statsd-mem-stats", opts.StatsdMemStats, "toggle sending memory and GC stats to statsd")
	flagSet.String("statsd-prefix", opts.StatsdPrefix, "prefix used for keys sent to statsd (%s for host replacement)")

	// End to end percentile flags
	e2eProcessingLatencyPercentiles := utils.FloatArray{}
	flagSet.Var(&e2eProcessingLatencyPercentiles, "e2e-processing-latency-percentile", "message processing time percentiles (as float (0, 1.0]) to track (can be specified multiple times or comma separated '1.0,0.99,0.95', default none)")
	flagSet.Duration("e2e-processing-latency-window-time", opts.E2EProcessingLatencyWindowTime, "calculate end to end latency quantiles for this duration of time (ie: 60s would only show quantile calculations from the past 60 seconds)")

	// TLS config
	flagSet.String("tls-cert", opts.TLSCert, "path to certificate file")
	flagSet.String("tls-key", opts.TLSKey, "path to key file")
	flagSet.String("tls-client-auth-policy", opts.TLSClientAuthPolicy, "client certificate auth policy ('require' or 'require-verify')")
	flagSet.String("tls-root-ca-file", opts.TLSRootCAFile, "path to certificate authority file")
	tlsRequired := tlsRequiredOption(opts.TLSRequired)
	tlsMinVersion := tlsMinVersionOption(opts.TLSMinVersion)
	flagSet.Var(&tlsRequired, "tls-required", "require TLS for client connections (true, false, tcp-https)")
	flagSet.Var(&tlsMinVersion, "tls-min-version", "minimum SSL/TLS version acceptable ('ssl3.0', 'tls1.0', 'tls1.1', or 'tls1.2')")

	// compression
	flagSet.Bool("deflate", opts.DeflateEnabled, "enable deflate feature negotiation (client compression)")
	flagSet.Int("max-deflate-level", opts.MaxDeflateLevel, "max deflate compression level a client can negotiate (> values == > eventbus CPU usage)")
	flagSet.Bool("snappy", opts.SnappyEnabled, "enable snappy feature negotiation (client compression)")

	return flagSet
}

type KVConfig map[string]interface{}

// Validate settings in the config file, and fatal on errors
func (cfg KVConfig) Validate() {
	// special validation/translation
	if v, exists := cfg["nsq_tls_required"]; exists {
		var t tlsRequiredOption
		err := t.Set(fmt.Sprintf("%v", v))
		if err == nil {
			cfg["nsq_tls_required"] = t.String()
		} else {
			log.Fatalf("ERROR: failed parsing tls required %v", v)
		}
	}
	if v, exists := cfg["nsq_tls_min_version"]; exists {
		var t tlsMinVersionOption
		err := t.Set(fmt.Sprintf("%v", v))
		if err == nil {
			newVal := fmt.Sprintf("%v", t.Get())
			if newVal != "0" {
				cfg["nsq_tls_min_version"] = newVal
			} else {
				delete(cfg, "nsq_tls_min_version")
			}
		} else {
			log.Fatalf("ERROR: failed parsing tls min version %v", v)
		}
	}
}
