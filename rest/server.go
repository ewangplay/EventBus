package rest

import (
	_ "net/http"

	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
	"github.com/fvbock/endless"
)

func StartServer(opts *config.EB_Options, logger i.ILogger, idcounter i.IIDCounter, messager i.IProducer, jobmgr i.IJobManager) (err error) {
	// Start Pprof debug server...
	if opts.PProfEnable {
		go func() {
			logger.Info("PPROF Debug Server Starting..., Listening on %s", opts.PProfAddress)
			err = endless.ListenAndServe(opts.PProfAddress, nil)
			if err != nil {
				logger.Error("PPROF Debug Server Return: %v", err)
			}
		}()
	}

	// Build REST router
	base_handler := NewBaseHandler(opts, logger, idcounter, messager, jobmgr)
	routes := map[string]i.Handler{
		"v1/event": &EventHandler{BaseHandler: base_handler},
		"v1/test":  &TestHandler{BaseHandler: base_handler},
	}
	router := NewRouter(opts, routes)

	// Start REST server...
	logger.Info("REST Server Start..., Listening on %s", opts.EBHTTPAddress)
	err = endless.ListenAndServe(opts.EBHTTPAddress, router)
	if err != nil {
		logger.Error("REST Server Return: %v", err)
		return err
	}

	return nil
}
