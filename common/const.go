package common

// Common constants define
const (
	LogID      = "logID"
	TimeCost   = "cost"
	RequestURL = "request_url"
	EventID    = "eventID"
	CreateTime = "createTime"
	ErrorCode  = "error_code"
	ErrorMsg   = "error_message"
	TimeFormat = "2006-01-02 15:04:05"
)

// Error message define
const (
	RequestBodyError = "[ERROR_INFO]Request body incorrect"
)

// Event status define
const (
	EventStatusUnset = iota
	EventStatusInit
	EventStatusDeal
	EventStatusSucc
	EventStatusFail
	EventStatusDrop
)

//Event retry policy define
const (
	CountRetryPolicy = 1 + iota
	ExpiredRetryPolicy
)
