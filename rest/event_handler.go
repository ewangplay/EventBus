package rest

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	c "github.com/ewangplay/eventbus/common"
)

// EventHandler structure
type EventHandler struct {
	*BaseHandler
}

// ProcessFunc http request callback function
func (eh *EventHandler) ProcessFunc(method string, resources []string, params map[string]string, body []byte, result map[string]interface{}) error {

	switch method {
	case "GET":
		if len(resources) < 3 {
			return fmt.Errorf("request uri incorrect")
		}
		eventID := resources[2]
		return eh.GetEvent(eventID, result)

	case "POST":
		return eh.AddEvent(body, result)
	}

	return fmt.Errorf("unsupported http method: %s", method)
}

// GetEvent return the specified event state
func (eh *EventHandler) GetEvent(eventID string, result map[string]interface{}) error {
	data, err := eh.Get(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "nil returned") {
			eh.Error("The event[%s] does not exist", eventID)
			return fmt.Errorf("[ERROR_INFO]event does not exist")
		}
		eh.Error("Get event[%s] from job manager error: %v", eventID, err)
		return err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		eh.Error("Parse event data[%s] error：%v", data, err)
		return err
	}

	return nil
}

// AddEvent publish event to eventbus
func (eh *EventHandler) AddEvent(body []byte, result map[string]interface{}) error {
	var err error
	logID := result[c.LogID]

	eh.Debug("[LogID:%v] Request body: %v", logID, string(body))

	//Parse event from request body
	var event c.EBEvent
	err = json.Unmarshal(body, &event)
	if err != nil {
		eh.Error("[LogID:%v] Parse request body error：%v", logID, err)
		return fmt.Errorf(c.RequestBodyError)
	}

	//Check and set the event's fields
	err = eh.checkAndSet(&event)
	if err != nil {
		eh.Error("[LogID:%v] Check request body error：%v", logID, err)
		return fmt.Errorf(c.RequestBodyError)
	}

	//Generate new id for eh event
	event.EventID, err = eh.NewEventID()
	if err != nil {
		eh.Error("[LogID:%v] New event id error：%v", logID, err)
		return err
	}

	//Register a job for eh event
	event.Status = c.EventStatusInit
	err = eh.Set(&event)
	if err != nil {
		eh.Error("[LogID:%v] Register event[%s:%s] error：%v", logID, event.Type, event.EventID, err)
		return err
	}

	//Publish eh event to MQ
	//topic := fmt.Sprintf("%s.%s", event.GetType(), event.GetSubject())
	err = eh.Publish(&c.EBMessage{Subject: event.GetType(), Data: event.GetData()})
	if err != nil {
		eh.Error("[LogID:%v] Publish event[%s:%s] error：%v", logID, event.Type, event.EventID, err)
		return err
	}

	//Response result
	result[c.EventID] = event.EventID
	result[c.CreateTime] = event.CreateTime

	return nil
}

func (eh *EventHandler) checkAndSet(event *c.EBEvent) error {
	// Step 1: Check
	if event.Subject == "" {
		return fmt.Errorf("request body must contain subject field")
	}

	if event.Type == "" {
		return fmt.Errorf("request body must contain type field")
	}

	if event.Body == nil {
		return fmt.Errorf("request body must contain body field")
	}

	// Step 2: Set
	event.CreateTime = time.Now().Format(c.TimeFormat)
	event.UpdateTime = event.CreateTime

	if event.RetryCount == 0 {
		event.RetryCount = eh.Opts.EBMaxRetryCount
	}
	if event.RetryInterval == 0 {
		event.RetryInterval = int64(eh.Opts.EBRetryInterval / time.Second)
	}
	if event.RetryTimeout == 0 {
		event.RetryTimeout = int64(eh.Opts.EBMaxRetryTimeout / time.Second)
	}
	if event.RetryPolicy == 0 {
		event.RetryPolicy = eh.Opts.EBRetryPolicy
	}

	return nil
}
