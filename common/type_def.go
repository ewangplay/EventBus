package common

import (
	"encoding/json"
	"fmt"
	"time"
)

//=========================================================================
//Event struct define
type EB_Event struct {
	Subject       string      `json:"subject"`
	EventId       string      `json:"event_id"`
	SessionId     string      `json:"session_id"`
	Type          string      `json:"type"`
	Status        int         `json:"status"`
	CreateTime    string      `json:"create_time"`
	UpdateTime    string      `json:"update_time"`
	RetryCount    int         `json:"retry_count"`
	RetryInterval int64       `json:"retry_interval"`
	RetryTimeout  int64       `json:"retry_timeout"`
	RetryPolicy   int         `json:"retry_policy"`
	TargetUrl     string      `json:"target_url"`
	Body          interface{} `json:"body"`
}

func (e EB_Event) String() string {
	return fmt.Sprintf("%s", e.GetData())
}

func (e EB_Event) GetId() string {
	return e.EventId
}

func (e EB_Event) GetSubject() string {
	return e.Subject
}

func (e EB_Event) GetType() string {
	return e.Type
}

func (e EB_Event) GetData() []byte {
	data, _ := json.Marshal(e)
	return data
}

func (e EB_Event) GetStatus() int {
	return e.Status
}

func (e EB_Event) GetRetryCount() int {
	return e.RetryCount
}

func (e EB_Event) GetRetryPolicy() int {
	return e.RetryPolicy
}

func (e EB_Event) GetRetryTimeout() int64 {
	return e.RetryTimeout
}

func (e EB_Event) GetRetryInterval() int64 {
	return e.RetryInterval
}

func (e EB_Event) GetCreateTime() int64 {
	if e.CreateTime == "" {
		e.CreateTime = time.Now().Format(TIME_FORMAT)
	}
	t, _ := time.Parse(TIME_FORMAT, e.CreateTime)
	return t.Unix()
}

func (e EB_Event) GetUpdateTime() int64 {
	if e.UpdateTime == "" {
		e.UpdateTime = time.Now().Format(TIME_FORMAT)
	}
	t, _ := time.Parse(TIME_FORMAT, e.UpdateTime)
	return t.Unix()
}

//=========================================================================
//Message struct define
type EB_Message struct {
	Subject string
	Data    []byte
}

func (m *EB_Message) GetSubject() string {
	return m.Subject
}

func (m *EB_Message) GetData() []byte {
	return m.Data
}
