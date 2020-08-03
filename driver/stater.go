package driver

import (
	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
)

// Stater ...
type Stater struct {
	i.Logger
	i.JobManager
	states chan i.Event
}

// NewStater new event stater instance
func NewStater(logger i.Logger, jobmgr i.JobManager) (*Stater, error) {
	s := &Stater{Logger: logger, JobManager: jobmgr}
	s.states = make(chan i.Event)

	go s.work()

	return s, nil
}

// Close closes current event stater instance
func (s *Stater) Close() error {
	if s.states != nil {
		close(s.states)
	}
	return nil
}

// PutState push one event into stater
func (s *Stater) PutState(event i.Event) {
	s.states <- event
}

func (s *Stater) work() {
	for state := range s.states {
		switch state.GetStatus() {
		case c.EventStatusDeal:
			fallthrough
		case c.EventStatusSucc:
			fallthrough
		case c.EventStatusDrop:
			s.Debug("Set event: %s", state.GetID())
			s.Set(state)
		case c.EventStatusFail:
			s.Debug("Fail event: %s", state.GetID())
			s.Fail(state)
		default:
			s.Error("message[%s] status incorrect!", state)
		}
	}
}
