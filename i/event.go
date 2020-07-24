package i

type IEvent interface {
	IMessage
	GetId() string
	GetType() string
	GetStatus() int
	GetRetryCount() int
	GetRetryPolicy() int
	GetRetryTimeout() int64
	GetRetryInterval() int64
	GetCreateTime() int64
	GetUpdateTime() int64
}
