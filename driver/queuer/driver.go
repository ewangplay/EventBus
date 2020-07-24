package queuer

import (
	"encoding/json"
	"fmt"
	"time"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
)

type Queuer struct {
	i.ILogger
	i.IStater
	i.IProducer
}

func NewQueuer(logger i.ILogger, stater i.IStater, producer i.IProducer) (*Queuer, error) {
	this := &Queuer{ILogger: logger, IStater: stater, IProducer: producer}
	return this, nil
}

func (this *Queuer) GetType() string {
	return "queuer"
}

func (this *Queuer) Process(messages <-chan i.IMessage) error {
	go func() {
		var err error
		var msg i.IMessage

		for msg = range messages {
			err = this.processMessage(msg)
			if err != nil {
				this.Error("Process message[%s:%s] error: %v", msg.GetSubject(), msg.GetData())
			}
		}
	}()

	return nil
}

func (this *Queuer) processMessage(msg i.IMessage) error {
	var err error

	if msg.GetSubject() != this.GetType() {
		this.Error("message type[%s] dismatch driver type[%s], drop it...", msg.GetSubject(), this.GetType())
		return fmt.Errorf("message type dismatch driver type")
	}

	//Parse event from message body
	var event c.EB_Event
	data := msg.GetData()
	err = json.Unmarshal(data, &event)
	if err != nil {
		this.Error("Parse message body[%v] error：%v", data, err)
		return err
	}

	if event.Status != c.ES_INIT {
		this.Error("invalid event state: %v", event.Status)
		return fmt.Errorf("invalid event state: %v", event.Status)
	}

	//Check the event status
	event.Status = c.ES_DEAL
	event.UpdateTime = time.Now().Format(c.TIME_FORMAT)
	this.PutState(event)

	this.Debug("Processing event [%s] start...", event.EventId)

	//Publish the event body into message queue
	data, err = json.Marshal(event.Body)
	if err != nil {
		this.Error("Marshal event body[%s] error: %v", event.Body, err)
		goto END
	}

	err = this.Publish(&c.EB_Message{Subject: event.Subject, Data: data})
	if err != nil {
		this.Error("Publish message [%s: %s] to queue error：%v", event.Subject, data, err)
		goto END
	}

	this.Debug("Processing event [%s] done", event.EventId)

END:
	//Set the event status
	if err == nil {
		event.Status = c.ES_SUCC
	} else {
		event.RetryCount--
		event.Status = c.ES_FAIL
	}
	event.UpdateTime = time.Now().Format(c.TIME_FORMAT)
	this.PutState(event)

	return err
}
