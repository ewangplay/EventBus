package rest

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	c "github.com/ewangplay/eventbus/common"
)

type EventHandler struct {
	*BaseHandler
}

func (this *EventHandler) ProcessFunc(method string, resources []string, params map[string]string, body []byte, result map[string]interface{}) error {

	switch method {
	case "GET":
		if len(resources) < 3 {
			return fmt.Errorf("request uri incorrect")
		}
		event_id := resources[2]
		return this.GetEvent(event_id, result)

	case "POST":
		return this.AddEvent(body, result)
	}

	return fmt.Errorf("unsupported http method: %s", method)
}

func (this *EventHandler) GetEvent(event_id string, result map[string]interface{}) error {
	data, err := this.Get(event_id)
	if err != nil {
		if strings.Contains(err.Error(), "nil returned") {
			this.Error("The event[%s] does not exist", event_id)
			return fmt.Errorf("[ERROR_INFO]event does not exist")
		} else {
			this.Error("Get event[%s] from job manager error: %v", event_id, err)
			return err
		}
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		this.Error("Parse event data[%s] error：%v", data, err)
		return err
	}

	return nil
}

func (this *EventHandler) AddEvent(body []byte, result map[string]interface{}) error {
	var err error
	log_id := result[c.LOG_ID]

	this.Debug("[LOG_ID:%v] Request body: %v", log_id, string(body))

	//Parse event from request body
	var event c.EB_Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		this.Error("[LOG_ID:%v] Parse request body error：%v", log_id, err)
		return fmt.Errorf(c.ERROR_BODY_ERR)
	}

	//Check and set the event's fields
	err = this.CheckAndSet(&event)
	if err != nil {
		this.Error("[LOG_ID:%v] Check request body error：%v", log_id, err)
		return fmt.Errorf(c.ERROR_BODY_ERR)
	}

	//Generate new id for this event
	event.EventId, err = this.NewEventId()
	if err != nil {
		this.Error("[LOG_ID:%v] New event id error：%v", log_id, err)
		return err
	}

	//Register a job for this event
	event.Status = c.ES_INIT
	err = this.Set(&event)
	if err != nil {
		this.Error("[LOG_ID:%v] Register event[%s:%s] error：%v", log_id, event.Type, event.EventId, err)
		return err
	}

	//Publish this event to MQ
	//topic := fmt.Sprintf("%s.%s", event.GetType(), event.GetSubject())
	err = this.Publish(&c.EB_Message{Subject: event.GetType(), Data: event.GetData()})
	if err != nil {
		this.Error("[LOG_ID:%v] Publish event[%s:%s] error：%v", log_id, event.Type, event.EventId, err)
		return err
	}

	//Response result
	result[c.EVENT_ID] = event.EventId
	result[c.CREATE_TIME] = event.CreateTime

	return nil
}

func (this *EventHandler) CheckAndSet(event *c.EB_Event) error {
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
	event.CreateTime = time.Now().Format(c.TIME_FORMAT)
	event.UpdateTime = event.CreateTime

	if event.RetryCount == 0 {
		event.RetryCount = this.Opts.EBMaxRetryCount
	}
	if event.RetryInterval == 0 {
		event.RetryInterval = int64(this.Opts.EBRetryInterval / time.Second)
	}
	if event.RetryTimeout == 0 {
		event.RetryTimeout = int64(this.Opts.EBMaxRetryTimeout / time.Second)
	}
	if event.RetryPolicy == 0 {
		event.RetryPolicy = this.Opts.EBRetryPolicy
	}

	return nil
}
