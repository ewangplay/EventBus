package i

// Event interface define
type Event interface {
	Message
	GetID() string
	GetType() string
	GetStatus() int
	GetRetryCount() int
	GetRetryPolicy() int
	GetRetryTimeout() int64
	GetRetryInterval() int64
	GetCreateTime() int64
	GetUpdateTime() int64
}
