package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/ewangplay/eventbus/adapter"
	msger "github.com/ewangplay/eventbus/adapter/messager"
	cfg "github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/driver"
	"github.com/ewangplay/eventbus/log"
	"github.com/ewangplay/eventbus/metadata"
	"github.com/ewangplay/eventbus/rest"
	options "github.com/mreiferson/go-options"
)

var (
	logger     *log.Logger
	idcounter  *adapter.Counter
	messager   *msger.Messager
	jobmgr     *adapter.JobManager
	dispatcher *driver.Dispatcher
)

func main() {
	var err error

	//New default options
	opts := cfg.NewOptions()
	opts.ServiceName = "eventbus"

	//New command line flag set
	flagSet := cfg.NewFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(metadata.GetVersionInfo("eventbus"))
		os.Exit(0)
	}

	//New config instance
	configFile := flagSet.Lookup("config").Value.String()
	config, err := cfg.New(configFile)
	if err != nil {
		fmt.Printf("[ERROR] Parse Configure File[%s] Error: %v\n", configFile, err)
		os.Exit(1)
	}

	//Resolves configuration values set via command line flags,
	//config files, and default struct values
	options.Resolve(opts, flagSet, config.GetConfig())

	fmt.Printf("EventBus Options: %+v\n", opts)

	//New logger instance
	logger, err = log.New(opts)
	if err != nil {
		fmt.Printf("[ERROR] Create logger Error: %v\n", err)
		goto END
	}

	opts.Logger = logger

	//New id counter instance
	idcounter, err = adapter.NewCounter(opts, logger)
	if err != nil {
		logger.Error("Create Counter Error: %v", err)
		goto END
	}

	//New messager instance
	messager, err = msger.NewMessager(opts, logger)
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

	if logger != nil {
		logger.Close()
	}
}
