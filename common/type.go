package common

import (
	"encoding/json"
	"fmt"
	"time"
)

//EBEvent struct define
type EBEvent struct {
	Subject       string      `json:"subject"`
	EventID       string      `json:"eventID"`
	Type          string      `json:"type"`
	Status        int         `json:"status"`
	CreateTime    string      `json:"createTime"`
	UpdateTime    string      `json:"updateTime"`
	RetryCount    int         `json:"retryCount"`
	RetryInterval int64       `json:"retryInterval"`
	RetryTimeout  int64       `json:"retryTimeout"`
	RetryPolicy   int         `json:"retryPolicy"`
	TargetURL     string      `json:"target_url"`
	Body          interface{} `json:"body"`
}

func (e EBEvent) String() string {
	return fmt.Sprintf("%s", e.GetData())
}

// GetID ...
func (e EBEvent) GetID() string {
	return e.EventID
}

// GetSubject ...
func (e EBEvent) GetSubject() string {
	return e.Subject
}

// GetType ...
func (e EBEvent) GetType() string {
	return e.Type
}

// GetData ...
func (e EBEvent) GetData() []byte {
	data, _ := json.Marshal(e)
	return data
}

// GetStatus ...
func (e EBEvent) GetStatus() int {
	return e.Status
}

// GetRetryCount ...
func (e EBEvent) GetRetryCount() int {
	return e.RetryCount
}

// GetRetryPolicy ...
func (e EBEvent) GetRetryPolicy() int {
	return e.RetryPolicy
}

// GetRetryTimeout ...
func (e EBEvent) GetRetryTimeout() int64 {
	return e.RetryTimeout
}

// GetRetryInterval ...
func (e EBEvent) GetRetryInterval() int64 {
	return e.RetryInterval
}

// GetCreateTime ...
func (e EBEvent) GetCreateTime() int64 {
	if e.CreateTime == "" {
		e.CreateTime = time.Now().Format(TimeFormat)
	}
	t, _ := time.Parse(TimeFormat, e.CreateTime)
	return t.Unix()
}

// GetUpdateTime ...
func (e EBEvent) GetUpdateTime() int64 {
	if e.UpdateTime == "" {
		e.UpdateTime = time.Now().Format(TimeFormat)
	}
	t, _ := time.Parse(TimeFormat, e.UpdateTime)
	return t.Unix()
}

//EBMessage struct define
type EBMessage struct {
	Subject string
	Data    []byte
}

// GetSubject ...
func (m *EBMessage) GetSubject() string {
	return m.Subject
}

// GetData ...
func (m *EBMessage) GetData() []byte {
	return m.Data
}
