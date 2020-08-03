package queuer

import (
	"encoding/json"
	"fmt"
	"time"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
)

// Queuer struct define
type Queuer struct {
	i.Logger
	i.Stater
	i.Producer
}

// NewQueuer ...
func NewQueuer(logger i.Logger, stater i.Stater, producer i.Producer) (*Queuer, error) {
	q := &Queuer{Logger: logger, Stater: stater, Producer: producer}
	return q, nil
}

// GetType ...
func (q *Queuer) GetType() string {
	return "queuer"
}

// Process ...
func (q *Queuer) Process(messages <-chan i.Message) error {
	go func() {
		var err error
		var msg i.Message

		for msg = range messages {
			err = q.processMessage(msg)
			if err != nil {
				q.Error("Process message[%s:%s] error: %v", msg.GetSubject(), msg.GetData())
			}
		}
	}()

	return nil
}

func (q *Queuer) processMessage(msg i.Message) error {
	var err error

	if msg.GetSubject() != q.GetType() {
		q.Error("message type[%s] dismatch driver type[%s], drop it...", msg.GetSubject(), q.GetType())
		return fmt.Errorf("message type dismatch driver type")
	}

	//Parse event from message body
	var event c.EBEvent
	data := msg.GetData()
	err = json.Unmarshal(data, &event)
	if err != nil {
		q.Error("Parse message body[%v] error：%v", data, err)
		return err
	}

	if event.Status != c.EventStatusInit {
		q.Error("invalid event state: %v", event.Status)
		return fmt.Errorf("invalid event state: %v", event.Status)
	}

	//Check the event status
	event.Status = c.EventStatusDeal
	event.UpdateTime = time.Now().Format(c.TimeFormat)
	q.PutState(event)

	q.Debug("Processing event [%s] start...", event.EventID)

	//Publish the event body into message queue
	data, err = json.Marshal(event.Body)
	if err != nil {
		q.Error("Marshal event body[%s] error: %v", event.Body, err)
		goto END
	}

	err = q.Publish(&c.EBMessage{Subject: event.Subject, Data: data})
	if err != nil {
		q.Error("Publish message [%s: %s] to queue error：%v", event.Subject, data, err)
		goto END
	}

	q.Debug("Processing event [%s] done", event.EventID)

END:
	//Set the event status
	if err == nil {
		event.Status = c.EventStatusSucc
	} else {
		event.RetryCount--
		event.Status = c.EventStatusFail
	}
	event.UpdateTime = time.Now().Format(c.TimeFormat)
	q.PutState(event)

	return err
}
