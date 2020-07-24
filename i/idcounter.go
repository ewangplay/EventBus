package i

type IIDCounter interface {
	NewEventId() (string, error)
}
