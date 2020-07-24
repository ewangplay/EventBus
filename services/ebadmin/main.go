package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/config"
	logger "github.com/ewangplay/eventbus/log"
	"github.com/ewangplay/eventbus/metadata"
	"github.com/ewangplay/eventbus/utils"
	options "github.com/mreiferson/go-options"
	"github.com/nsqio/nsq/nsqadmin"
)

var (
	flagSet = flag.NewFlagSet("nsqadmin", flag.ExitOnError)

	configFile  = flagSet.String("config", "", "path to config file")
	showVersion = flagSet.Bool("version", false, "print version string")

	logLevel  = flagSet.String("log-level", "info", "set log verbosity: debug, info, warn, error, or fatal")
	logPrefix = flagSet.String("log-prefix", "[eventbus-admin] ", "log message prefix")

	httpAddress = flagSet.String("http-address", "0.0.0.0:4171", "<addr>:<port> to listen on for HTTP clients")

	graphiteURL   = flagSet.String("graphite-url", "", "graphite HTTP address")
	proxyGraphite = flagSet.Bool("proxy-graphite", false, "proxy HTTP requests to graphite")

	statsdCounterFormat = flagSet.String("statsd-counter-format", "stats.counters.%s.count", "The counter stats key formatting applied by the implementation of statsd. If no formatting is desired, set this to an empty string.")
	statsdGaugeFormat   = flagSet.String("statsd-gauge-format", "stats.gauges.%s", "The gauge stats key formatting applied by the implementation of statsd. If no formatting is desired, set this to an empty string.")
	statsdPrefix        = flagSet.String("statsd-prefix", "nsq.%s", "prefix used for keys sent to statsd (%s for host replacement, must match nsqd)")
	statsdInterval      = flagSet.Duration("statsd-interval", 60*time.Second, "time interval nsqd is configured to push to statsd (must match nsqd)")

	notificationHTTPEndpoint = flagSet.String("notification-http-endpoint", "", "HTTP endpoint (fully qualified) to which POST notifications of admin actions will be sent")

	httpConnectTimeout = flagSet.Duration("http-client-connect-timeout", 2*time.Second, "timeout for HTTP connect")
	httpRequestTimeout = flagSet.Duration("http-client-request-timeout", 5*time.Second, "timeout for HTTP request")

	httpClientTLSInsecureSkipVerify = flagSet.Bool("http-client-tls-insecure-skip-verify", false, "configure the HTTP client to skip verification of TLS certificates")
	httpClientTLSRootCAFile         = flagSet.String("http-client-tls-root-ca-file", "", "path to CA file for the HTTP client")
	httpClientTLSCert               = flagSet.String("http-client-tls-cert", "", "path to certificate file for the HTTP client")
	httpClientTLSKey                = flagSet.String("http-client-tls-key", "", "path to key file for the HTTP client")

	allowConfigFromCIDR = flagSet.String("allow-config-from-cidr", "127.0.0.1/8", "A CIDR from which to allow HTTP requests to the /config endpoint")
	aclHttpHeader       = flagSet.String("acl-http-header", "X-Forwarded-User", "HTTP header to check for authenticated admin users")

	adminUsers              = utils.StringArray{}
	nsqlookupdHTTPAddresses = utils.StringArray{}
	nsqdHTTPAddresses       = utils.StringArray{}
)

func init() {
	flagSet.Var(&nsqlookupdHTTPAddresses, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
	flagSet.Var(&nsqdHTTPAddresses, "nsqd-http-address", "nsqd HTTP address (may be given multiple times)")
	flagSet.Var(&adminUsers, "admin-user", "admin user (may be given multiple times; if specified, only these users will be able to perform privileged actions; acl-http-header is used to determine the authenticated user)")
}

func main() {
	flagSet.Parse(os.Args[1:])

	if *showVersion {
		fmt.Println(metadata.GetVersionInfo(common.ProgramNameEBAdmin))
		return
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var cfg map[string]interface{}
	if *configFile != "" {
		_, err := toml.DecodeFile(*configFile, &cfg)
		if err != nil {
			log.Fatalf("ERROR: failed to load config file %s - %s", *configFile, err)
		}
	}

	opts := nsqadmin.NewOptions()
	options.Resolve(opts, flagSet, cfg)

	ebopts := config.NewOptions()
	ebopts.ServiceName = "eventbus-admin"
	ebopts.LogPrefix = opts.LogPrefix
	ebopts.LogLevel = config.LogLevel(opts.LogLevel)
	logPath, ok := cfg["log-path"]
	if ok {
		slogPath, ok := logPath.(string)
		if ok && slogPath != "" {
			ebopts.LogPath = slogPath
			log.Printf("LogPath: %s", ebopts.LogPath)
		}
	}
	logMaxSize, ok := cfg["log-max-size"]
	if ok {
		nlogMaxSize, ok := logMaxSize.(int64)
		if ok && nlogMaxSize != 0 {
			ebopts.LogMaxSize = nlogMaxSize
			log.Printf("LogMaxSize: %d", ebopts.LogMaxSize)
		}
	}
	logRotateDaily, ok := cfg["log-rotate-daily"]
	if ok {
		blogRotateDaily, ok := logRotateDaily.(bool)
		if ok {
			ebopts.LogRotateDaily = blogRotateDaily
			log.Printf("LogRotateDaily: %v", ebopts.LogRotateDaily)
		}
	}
	logMaxAge, ok := cfg["log-max-age"]
	if ok {
		nlogMaxAge, ok := logMaxAge.(int)
		if ok && nlogMaxAge != 0 {
			ebopts.LogMaxAge = nlogMaxAge
			log.Printf("LogMaxAge: %d", ebopts.LogMaxAge)
		}
	}

	logger, err := logger.New(ebopts)
	if err == nil {
		opts.Logger = logger
	}

	nsqadmin, err := nsqadmin.New(opts)
	if err != nil {
		log.Fatalf("ERROR: new nsqdadmin instance fail: %v", err)
	}
	nsqadmin.Main()
	<-exitChan
	nsqadmin.Exit()
}