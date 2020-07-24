package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/ewangplay/eventbus/adapter"
	"github.com/ewangplay/eventbus/adapter/nsqd"
	"github.com/ewangplay/eventbus/common"
	cfg "github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/driver"
	"github.com/ewangplay/eventbus/log"
	"github.com/ewangplay/eventbus/metadata"
	"github.com/ewangplay/eventbus/rest"
	options "github.com/mreiferson/go-options"
)

var (
	logger     *log.Logger
	nsqdSrv    *nsqd.NsqdService
	idcounter  *adapter.IDCounter
	messager   *adapter.Messager
	jobmgr     *adapter.JobManager
	dispatcher *driver.Dispatcher
)

func main() {
	var err error

	//New default options
	opts := cfg.NewOptions()
	opts.ServiceName = "eventbus-node"

	//New command line flag set
	flagSet := cfg.NewFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(metadata.GetVersionInfo(common.ProgramNameEBNode))
		os.Exit(0)
	}

	//New config instance
	configFile := flagSet.Lookup("config").Value.String()
	config, err := cfg.New(configFile)
	if err != nil {
		fmt.Printf("[ERROR] Parse Configure File[%s] Error: %v\n", configFile, err)
		os.Exit(1)
	}

	var kvcfg cfg.KVConfig
	kvcfg = config.GetConfig()
	fmt.Printf("kvcfg: %v\n", kvcfg)
	kvcfg.Validate()

	//Resolves configuration values set via command line flags,
	//config files, and default struct values
	options.Resolve(opts, flagSet, kvcfg)

	fmt.Printf("EventBus Options: %+v\n", opts)

	//New logger instance
	logger, err = log.New(opts)
	if err != nil {
		fmt.Printf("[ERROR] Create logger Error: %v\n", err)
		goto END
	}

	opts.Logger = logger

	//New nsqd service
	if opts.NSQEnable {
		nsqdSrv, err = nsqd.NewNsqdService()
		if err != nil {
			logger.Error("Create NSQD Service Error: %v", err)
			goto END
		}
		err = nsqdSrv.Start(opts)
		if err != nil {
			logger.Error("Start NSQD Service Error: %v", err)
			goto END
		}
	}

	//New id counter instance
	idcounter, err = adapter.NewIdCounter(opts, logger)
	if err != nil {
		logger.Error("Create IDCounter Error: %v", err)
		goto END
	}

	//New messager instance
	messager, err = adapter.NewMessager(opts, logger)
	if err != nil {
		logger.Error("Create Messager Error: %v", err)
		goto END
	}

	//New job manager instance
	jobmgr, err = adapter.NewJobManager(opts, logger, messager)
	if err != nil {
		logger.Error("Create JobManager Error: %v", err)
		goto END
	}

	//New event dispatcher instance
	dispatcher, err = driver.NewDispatcher(opts, logger, messager, jobmgr)
	if err != nil {
		logger.Error("Create Event Dispatcher Error: %v", err)
		goto END
	}

	//Start REST Server
	rest.StartServer(opts, logger, idcounter, messager, jobmgr)

END:
	cleanup()
	if err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func cleanup() {
	if dispatcher != nil {
		dispatcher.Close()
	}

	if jobmgr != nil {
		jobmgr.Close()
	}

	if messager != nil {
		messager.Close()
	}

	if idcounter != nil {
		idcounter.Close()
	}

	if nsqdSrv != nil {
		nsqdSrv.Stop()
	}

	if logger != nil {
		logger.Close()
	}
}
