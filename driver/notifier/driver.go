package notifier

import (
	"encoding/json"
	"fmt"
	"time"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
	"github.com/ewangplay/eventbus/utils"
)

// Notifier struct define
type Notifier struct {
	i.Logger
	i.Stater
}

// NewNotifier new notifier driver
func NewNotifier(logger i.Logger, stater i.Stater) (*Notifier, error) {
	n := &Notifier{Logger: logger, Stater: stater}
	return n, nil
}

// GetType get driver's type
func (n *Notifier) GetType() string {
	return "notifier"
}

// Process ...
func (n *Notifier) Process(messages <-chan i.Message) error {
	go func() {
		var err error
		var msg i.Message

		for msg = range messages {
			err = n.processMessage(msg)
			if err != nil {
				n.Error("Process message[%s:%s] error: %v", msg.GetSubject(), msg.GetData())
			}
		}
	}()

	return nil
}

func (n *Notifier) processMessage(msg i.Message) error {
	var err error

	if msg.GetSubject() != n.GetType() {
		n.Error("message subject[%s] dismatch driver type[%s], drop it...", msg.GetSubject(), n.GetType())
		return fmt.Errorf("message subject dismatch driver type")
	}

	//Parse event from message body
	var event c.EBEvent
	data := msg.GetData()
	err = json.Unmarshal(data, &event)
	if err != nil {
		n.Error("Parse message body[%v] error：%v", data, err)
		return err
	}

	if event.Status != c.EventStatusInit {
		n.Error("invalid event status: %v", event.Status)
		return fmt.Errorf("invalid event status: %v", event.Status)
	}

	//Check the event status
	event.Status = c.EventStatusDeal
	event.UpdateTime = time.Now().Format(c.TimeFormat)
	n.PutState(event)

	n.Debug("Processing event [%s] start...", event.EventID)

	//PUsh the event body to target url
	var result []byte

	data, err = json.Marshal(event.Body)
	if err != nil {
		n.Error("Marshal event body[%v] error: %v", event.Body, err)
		goto END
	}

	result, err = utils.HTTPPost(event.TargetURL, data)
	if err != nil {
		n.Error("Post event to target url[%s] error: %v", event.TargetURL, err)
		goto END
	}
	if result != nil {
		n.Info("Response result: %s", string(result))
	}

	n.Debug("Processing event [%s] done", event.EventID)

END:
	//Set the event status
	if err == nil {
		event.Status = c.EventStatusSucc
	} else {
		event.RetryCount--
		event.Status = c.EventStatusFail
	}
	event.UpdateTime = time.Now().Format(c.TimeFormat)
	n.PutState(event)

	return err
}
