package i

type IStater interface {
	PutState(IEvent)
}
