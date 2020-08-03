package i

// JobManager interface define
type JobManager interface {
	Set(Event) error
	Fail(Event) error
	Get(string) ([]byte, error)
}
