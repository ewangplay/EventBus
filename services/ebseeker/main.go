package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/config"
	logger "github.com/ewangplay/eventbus/log"
	"github.com/ewangplay/eventbus/metadata"
	"github.com/judwhite/go-svc/svc"
	options "github.com/mreiferson/go-options"
	"github.com/nsqio/nsq/nsqlookupd"
)

func nsqlookupdFlagSet(opts *nsqlookupd.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("nsqlookupd", flag.ExitOnError)

	flagSet.String("config", "", "path to config file")
	flagSet.Bool("version", false, "print version string")

	flagSet.String("log-level", "info", "set log verbosity: debug, info, warn, error, or fatal")
	flagSet.String("log-prefix", "[eventbus-seeker] ", "log message prefix")

	flagSet.String("tcp-address", opts.TCPAddress, "<addr>:<port> to listen on for TCP clients")
	flagSet.String("http-address", opts.HTTPAddress, "<addr>:<port> to listen on for HTTP clients")
	flagSet.String("broadcast-address", opts.BroadcastAddress, "address of this lookupd node, (default to the OS hostname)")

	flagSet.Duration("inactive-producer-timeout", opts.InactiveProducerTimeout, "duration of time a producer will remain in the active list since its last ping")
	flagSet.Duration("tombstone-lifetime", opts.TombstoneLifetime, "duration of time a producer will remain tombstoned if registration remains")

	return flagSet
}

type program struct {
	nsqlookupd *nsqlookupd.NSQLookupd
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	opts := nsqlookupd.NewOptions()

	flagSet := nsqlookupdFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(metadata.GetVersionInfo(common.ProgramNameEBSeeker))
		os.Exit(0)
	}

	var cfg map[string]interface{}
	configFile := flagSet.Lookup("config").Value.String()
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &cfg)
		if err != nil {
			log.Fatalf("ERROR: failed to load config file %s - %s", configFile, err.Error())
		}
	}

	options.Resolve(opts, flagSet, cfg)

	ebopts := config.NewOptions()
	ebopts.ServiceName = "eventbus-seeker"
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

	daemon, err := nsqlookupd.New(opts)
	if err != nil {
		log.Fatalf("ERROR: new nsqlookupd install fail: %v", err)
	}
	daemon.Main()
	p.nsqlookupd = daemon
	return nil
}

func (p *program) Stop() error {
	if p.nsqlookupd != nil {
		p.nsqlookupd.Exit()
	}
	return nil
}
