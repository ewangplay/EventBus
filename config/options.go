package config

import (
	"crypto/tls"
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
	// basic options
	NSQEnable  bool `flag:"nsq-enable"`
	NSQCluster bool `flag:"nsq-cluster"`

	NSQMaxInFlight int `flag:"nsq-max-in-flight"`

	NSQTCPAddress          string   `flag:"nsq-tcp-address"`
	NSQHTTPAddress         string   `flag:"nsq-http-address"`
	NSQHTTPSAddress        string   `flag:"nsq-https-address"`
	NSQBroadcastAddress    string   `flag:"nsq-broadcast-address"`
	NSQLookupdTCPAddresses []string `flag:"nsq-lookupd-tcp-address" cfg:"nsq_lookupd_tcp_addresses"`
	NSQAuthHTTPAddresses   []string `flag:"nsq-auth-http-address" cfg:"nsq_auth_http_addresses"`

	HTTPClientConnectTimeout time.Duration `flag:"http-client-connect-timeout" cfg:"nsq_http_client_connect_timeout"`
	HTTPClientRequestTimeout time.Duration `flag:"http-client-request-timeout" cfg:"nsq_http_client_request_timeout"`

	// diskqueue options
	DataPath        string        `flag:"data-path" cfg:"nsq_data_path"`
	MemQueueSize    int64         `flag:"mem-queue-size" cfg:"nsq_mem_queue_size"`
	MaxBytesPerFile int64         `flag:"max-bytes-per-file" cfg:"naq_max_bytes_per_file"`
	SyncEvery       int64         `flag:"sync-every" cfg:"nsq_sync_every"`
	SyncTimeout     time.Duration `flag:"sync-timeout" cfg:"nsq_sync_timeout"`

	QueueScanInterval        time.Duration
	QueueScanRefreshInterval time.Duration
	QueueScanSelectionCount  int
	QueueScanWorkerPoolMax   int
	QueueScanDirtyPercent    float64

	// msg and command options
	MsgTimeout    time.Duration `flag:"msg-timeout" cfg:"nsq_msg_timeout"`
	MaxMsgTimeout time.Duration `flag:"max-msg-timeout" cfg:"nsq_max_msg_timeout"`
	MaxMsgSize    int64         `flag:"max-msg-size" cfg:"nsq_max_msg_size"`
	MaxBodySize   int64         `flag:"max-body-size" cfg:"nsq_max_body_size"`
	MaxReqTimeout time.Duration `flag:"max-req-timeout" cfg:"nsq_max_req_timeout"`
	ClientTimeout time.Duration

	// client overridable configuration options
	MaxHeartbeatInterval   time.Duration `flag:"max-heartbeat-interval" cfg:"nsq_max_heartbeat_interval"`
	MaxRdyCount            int64         `flag:"max-rdy-count" cfg:"nsq_max_rdy_count"`
	MaxOutputBufferSize    int64         `flag:"max-output-buffer-size" cfg:"nsq_max_output_buffer_size"`
	MaxOutputBufferTimeout time.Duration `flag:"max-output-buffer-timeout" cfg:"nsq_max_output_buffer_timeout"`

	// statsd integration
	StatsdAddress  string        `flag:"statsd-address" cfg:"nsq_statsd_address"`
	StatsdPrefix   string        `flag:"statsd-prefix" cfg:"nsq_statsd_prefix"`
	StatsdInterval time.Duration `flag:"statsd-interval" cfg:"nsq_statsd_interval"`
	StatsdMemStats bool          `flag:"statsd-mem-stats" cfg:"nsq_statsd_mem_stats"`

	// e2e message latency
	E2EProcessingLatencyWindowTime  time.Duration `flag:"e2e-processing-latency-window-time" cfg:"nsq_e2e_processing_latency_window_time"`
	E2EProcessingLatencyPercentiles []float64     `flag:"e2e-processing-latency-percentile" cfg:"nsq_e2e_processing_latency_percentiles"`

	// TLS config
	TLSCert             string `flag:"tls-cert" cfg:"nsq_tls_cert"`
	TLSKey              string `flag:"tls-key" cfg:"nsq_tls_key"`
	TLSClientAuthPolicy string `flag:"tls-client-auth-policy" cfg:"nsq_tls_client_auth_policy"`
	TLSRootCAFile       string `flag:"tls-root-ca-file" cfg:"nsq_tls_root_ca_file"`
	TLSRequired         int    `flag:"tls-required" cfg:"nsq_tls_required"`
	TLSMinVersion       uint16 `flag:"tls-min-version" cfg:"nsq_tls_min_version"`

	// compression
	DeflateEnabled  bool `flag:"deflate" cfg:"nsq_deflate"`
	MaxDeflateLevel int  `flag:"max-deflate-level" cfg:"nsq_max_deflate_level"`
	SnappyEnabled   bool `flag:"snappy" cfg:"nsq_snappy"`
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

		NSQTCPAddress:   "0.0.0.0:4150",
		NSQHTTPAddress:  "0.0.0.0:4151",
		NSQHTTPSAddress: "0.0.0.0:4152",

		NSQLookupdTCPAddresses: make([]string, 0),
		NSQAuthHTTPAddresses:   make([]string, 0),

		HTTPClientConnectTimeout: 2 * time.Second,
		HTTPClientRequestTimeout: 5 * time.Second,

		DataPath:        "/opt/ewangplay/eventbus/data",
		MemQueueSize:    10000,
		MaxBytesPerFile: 100 * 1024 * 1024,
		SyncEvery:       2500,
		SyncTimeout:     2 * time.Second,

		QueueScanInterval:        100 * time.Millisecond,
		QueueScanRefreshInterval: 5 * time.Second,
		QueueScanSelectionCount:  20,
		QueueScanWorkerPoolMax:   4,
		QueueScanDirtyPercent:    0.25,

		MsgTimeout:    60 * time.Second,
		MaxMsgTimeout: 15 * time.Minute,
		MaxMsgSize:    1024 * 1024,
		MaxBodySize:   5 * 1024 * 1024,
		MaxReqTimeout: 1 * time.Hour,
		ClientTimeout: 60 * time.Second,

		MaxHeartbeatInterval:   60 * time.Second,
		MaxRdyCount:            2500,
		MaxOutputBufferSize:    64 * 1024,
		MaxOutputBufferTimeout: 1 * time.Second,

		StatsdPrefix:   "nsq.%s",
		StatsdInterval: 60 * time.Second,
		StatsdMemStats: true,

		E2EProcessingLatencyWindowTime: time.Duration(10 * time.Minute),

		DeflateEnabled:  true,
		MaxDeflateLevel: 6,
		SnappyEnabled:   true,

		TLSMinVersion: tls.VersionTLS10,
	}
}
