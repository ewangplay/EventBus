package driver

import (
	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
)

type DriverStater struct {
	i.ILogger
	i.IJobManager
	states chan i.IEvent
}

func NewDriverStater(logger i.ILogger, jobmgr i.IJobManager) (*DriverStater, error) {
	this := &DriverStater{ILogger: logger, IJobManager: jobmgr}
	this.states = make(chan i.IEvent)

	go this.work()

	return this, nil
}

func (this *DriverStater) Close() error {
	if this.states != nil {
		close(this.states)
	}
	return nil
}

func (this *DriverStater) PutState(event i.IEvent) {
	this.states <- event
}

func (this *DriverStater) work() {
	for state := range this.states {
		switch state.GetStatus() {
		case c.ES_DEAL:
			fallthrough
		case c.ES_SUCC:
			fallthrough
		case c.ES_DROP:
			this.Debug("Set event: %s", state.GetId())
			this.Set(state)
		case c.ES_FAIL:
			this.Debug("Fail event: %s", state.GetId())
			this.Fail(state)
		default:
			this.Error("message[%s] status incorrect!", state)
		}
	}
}
