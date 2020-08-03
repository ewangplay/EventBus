package driver

import (
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

// Dispatcher ...
type Dispatcher struct {
	i.Logger
	i.Messager
	opts          *config.EBOptions
	driverFactory *Factory
}

// NewDispatcher new event dispatcher instance
func NewDispatcher(opts *config.EBOptions, logger i.Logger, messager i.Messager, jobmgr i.JobManager) (*Dispatcher, error) {
	//New driver factory instance
	driverFactory, err := NewFactory(logger, jobmgr, messager)
	if err != nil {
		logger.Error("Create Factory Error: %v", err)
		return nil, err
	}

	d := &Dispatcher{Logger: logger, Messager: messager, opts: opts}

	d.driverFactory = driverFactory

	go d.work()

	return d, nil
}

// Close closes current event dispatcher instance
func (d *Dispatcher) Close() error {
	d.Info("Driver dispatcher will close")
	return d.driverFactory.Close()
}

func (d *Dispatcher) work() error {
	for _, subject := range d.opts.Drivers {
		//subscribe the subject
		chMsg, err := d.Subscribe(subject)
		if err != nil {
			d.Error("subscribe message for subject[%v] error: %v", subject, err)
			return err
		}

		//create the subject driver
		driver, err := d.driverFactory.CreateDriver(subject)
		if err != nil {
			d.Error("create driver for subject[%v] error: %v", subject, err)
			return err
		}

		//process the messages using the driver
		err = driver.Process(chMsg)
		if err != nil {
			d.Error("driver[%v] start working error: %v", subject, err)
			return err
		}
	}

	return nil
}
