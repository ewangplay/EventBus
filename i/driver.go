package i

// Driver interface define
type Driver interface {
	GetType() string
	Process(messages <-chan Message) error
}
