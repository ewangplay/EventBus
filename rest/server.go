package rest

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
	"github.com/fvbock/endless"
)

// StartServer ...
func StartServer(opts *config.EBOptions, logger i.Logger, idcounter i.Counter, messager i.Producer, jobmgr i.JobManager) (err error) {
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
	baseHandler := NewBaseHandler(opts, logger, idcounter, messager, jobmgr)
	routes := map[string]i.Handler{
		"v1/event": &EventHandler{BaseHandler: baseHandler},
		"v1/test":  &TestHandler{BaseHandler: baseHandler},
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
