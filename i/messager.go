package i

type IMessage interface {
	GetSubject() string
	GetData() []byte
}

type IContext interface {
	CreateProducer() (IProducer, error)
	CreateConsumer(subject string) (IConsumer, error)
}

type IProducer interface {
	Publish(msg IMessage) error
	Close() error
}

type IConsumer interface {
	GetMessage() <-chan IMessage
	Close() error
}

type IMessager interface {
	Publish(msg IMessage) error
	Subscribe(subject string) (<-chan IMessage, error)
	Unsubscribe(subject string) error
	Close() error
}
