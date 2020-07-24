package driver

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

type Dispatcher struct {
	i.ILogger
	i.IMessager
	opts          *config.EB_Options
	driverFactory *DriverFactory
}

func NewDispatcher(opts *config.EB_Options, logger i.ILogger, messager i.IMessager, jobmgr i.IJobManager) (*Dispatcher, error) {
	//New driver factory instance
	driverFactory, err := NewDriverFactory(logger, jobmgr, messager)
	if err != nil {
		logger.Error("Create DriverFactory Error: %v", err)
		return nil, err
	}

	this := &Dispatcher{ILogger: logger, IMessager: messager, opts: opts}

	this.driverFactory = driverFactory

	go this.work()

	return this, nil
}

func (this *Dispatcher) Close() error {
	this.Info("Driver dispatcher will close")
	return this.driverFactory.Close()
}

func (this *Dispatcher) work() error {
	for _, subject := range this.opts.Drivers {
		//subscribe the subject
		chMsg, err := this.Subscribe(subject)
		if err != nil {
			this.Error("subscribe message for subject[%v] error: %v", subject, err)
			return err
		}

		//create the subject driver
		driver, err := this.driverFactory.CreateDriver(subject)
		if err != nil {
			this.Error("create driver for subject[%v] error: %v", subject, err)
			return err
		}

		//process the messages using the driver
		err = driver.Process(chMsg)
		if err != nil {
			this.Error("driver[%v] start working error: %v", subject, err)
			return err
		}
	}

	return nil
}
