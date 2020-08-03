package i

// Logger interface define
type Logger interface {
	Fatal(format string, args ...interface{}) error
	Error(format string, args ...interface{}) error
	Warn(format string, args ...interface{}) error
	Info(format string, args ...interface{}) error
	Debug(format string, args ...interface{}) error
	Output(maxdepth int, s string) error
}
