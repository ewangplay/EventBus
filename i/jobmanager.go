package i

type IJobManager interface {
	Set(IEvent) error
	Fail(IEvent) error
	Get(string) ([]byte, error)
}
