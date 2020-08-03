package i

// Message interface define
type Message interface {
	GetSubject() string
	GetData() []byte
}

// Context interface define
type Context interface {
	CreateProducer() (Producer, error)
	CreateConsumer(subject string) (Consumer, error)
}

// Producer interface define
type Producer interface {
	Publish(msg Message) error
	Close() error
}

// Consumer interface define
type Consumer interface {
	GetMessage() <-chan Message
	Close() error
}

// Messager interface define
type Messager interface {
	Publish(msg Message) error
	Subscribe(subject string) (<-chan Message, error)
	Unsubscribe(subject string) error
	Close() error
}
