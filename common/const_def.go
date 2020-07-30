package common

//constant define
const (
	LOG_ID      = "log_id"
	TIME_COST   = "cost"
	REQUEST_URL = "request_url"
	EVENT_ID    = "event_id"
	CREATE_TIME = "create_time"
	ERROR_CODE  = "error_code"
	ERROR_MSG   = "error_message"
)

//time format
const TIME_FORMAT = "2006-01-02 15:04:05"

//error define
const (
	ERROR_BODY_ERR = "[ERROR_INFO]Request body incorrect"
)

//Event status define
const (
	ES_UNSET = iota
	ES_INIT
	ES_DEAL
	ES_SUCC
	ES_FAIL
	ES_DROP
)

//Event retry policy define
const (
	ERP_COUNT = 1 + iota
	ERP_TIMEOUT
)