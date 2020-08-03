package driver

import (
	"fmt"

	"github.com/ewangplay/eventbus/driver/notifier"
	"github.com/ewangplay/eventbus/driver/queuer"
	"github.com/ewangplay/eventbus/i"
)

// Supported driver list
const (
	NotifierDriver = "notifier"
	QueuerDriver   = "queuer"
)

// Factory driver factory structure
type Factory struct {
	i.Logger
	i.Producer
	stater *Stater
}

// NewFactory new driver factory instance
func NewFactory(logger i.Logger, jobmgr i.JobManager, producer i.Producer) (*Factory, error) {

	stater, err := NewStater(logger, jobmgr)
	if err != nil {
		logger.Error("Create Stater Error: %v", err)
		return nil, err
	}

	factory := &Factory{Logger: logger, Producer: producer}
	factory.stater = stater

	return factory, nil
}

// Close the current driver factory instance
func (f *Factory) Close() error {
	return f.stater.Close()
}

// CreateDriver craate specified event driver using current driver factory
func (f *Factory) CreateDriver(subject string) (i.Driver, error) {
	switch subject {
	case NotifierDriver:
		return notifier.NewNotifier(f.Logger, f.stater)
	case QueuerDriver:
		return queuer.NewQueuer(f.Logger, f.stater, f.Producer)
	}

	f.Error("not supported driver: %s", subject)
	return nil, fmt.Errorf("not supported driver: %s", subject)
}
