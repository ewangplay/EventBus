package driver

import (
	"fmt"

	"github.com/ewangplay/eventbus/driver/notifier"
	"github.com/ewangplay/eventbus/driver/queuer"
	"github.com/ewangplay/eventbus/i"
)

const (
	DRIVER_NOTIFIER = "notifier"
	DRIVER_QUEUER   = "queuer"
)

type DriverFactory struct {
	i.ILogger
	i.IProducer
	stater *DriverStater
}

func NewDriverFactory(logger i.ILogger, jobmgr i.IJobManager, producer i.IProducer) (*DriverFactory, error) {

	stater, err := NewDriverStater(logger, jobmgr)
	if err != nil {
		logger.Error("Create DriverStater Error: %v", err)
		return nil, err
	}

	factory := &DriverFactory{ILogger: logger, IProducer: producer}
	factory.stater = stater

	return factory, nil
}

func (this *DriverFactory) Close() error {
	return this.stater.Close()
}

func (this *DriverFactory) CreateDriver(subject string) (i.IDriver, error) {
	switch subject {
	case DRIVER_NOTIFIER:
		return notifier.NewNotifier(this.ILogger, this.stater)
	case DRIVER_QUEUER:
		return queuer.NewQueuer(this.ILogger, this.stater, this.IProducer)
	}

	this.Error("not supported driver: %s", subject)
	return nil, fmt.Errorf("not supported driver: %s", subject)
}
