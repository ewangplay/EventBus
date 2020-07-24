package i

type IDriver interface {
	GetType() string
	Process(messages <-chan IMessage) error
}
