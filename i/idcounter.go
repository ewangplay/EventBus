package i

// Counter interface define
type Counter interface {
	NewEventID() (string, error)
}
