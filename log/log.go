package log

import (
	"fmt"
	"path/filepath"

	"github.com/arxanfintech/go-logger"
	"github.com/ewangplay/eventbus/config"
)

// Logger struct define, implement Logger interface
type Logger struct {
	serviceName  string
	cfgLevel     config.LogLevel
	logger       *logger.Logger
	rotateWriter *RotateWriter
}

// New Logger instance
func New(opts *config.EBOptions) (log *Logger, err error) {
	service := opts.ServiceName
	filename := fmt.Sprintf("%s.log", service)
	logFile := filepath.Join(opts.LogPath, filename)

	rotateWriter, err := NewRotateWriter(logFile,
		opts.LogMaxSize,
		opts.LogRotateDaily)
	if err != nil {
		return nil, err
	}

	color := 0
	logger, err := logger.New(service, color, rotateWriter)
	if err != nil {
		return nil, err
	}

	log = &Logger{
		serviceName:  service,
		cfgLevel:     config.ParseLogLevel(opts.LogLevel),
		logger:       logger,
		rotateWriter: rotateWriter,
	}

	return log, nil
}

// Close logger instance
func (l *Logger) Close() (err error) {
	return l.rotateWriter.Close()
}

func (l *Logger) log(level config.LogLevel, format string, args ...interface{}) (err error) {
	if level < l.cfgLevel {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	fmt.Println(msg)

	s := fmt.Sprintf("%s: %s", level, msg)

	l.Output(3, s)

	return nil
}

// Output ...
func (l *Logger) Output(maxdepth int, s string) error {
	maxdepth += 2
	l.logger.Log("", maxdepth, s)
	return nil
}

// Fatal ...
func (l *Logger) Fatal(format string, args ...interface{}) (err error) {
	return l.log(config.FATAL, format, args...)
}

// Error ...
func (l *Logger) Error(format string, args ...interface{}) (err error) {
	return l.log(config.ERROR, format, args...)
}

// Warn ...
func (l *Logger) Warn(format string, args ...interface{}) (err error) {
	return l.log(config.WARN, format, args...)
}

// Info ...
func (l *Logger) Info(format string, args ...interface{}) (err error) {
	return l.log(config.INFO, format, args...)
}

// Debug ...
func (l *Logger) Debug(format string, args ...interface{}) (err error) {
	return l.log(config.DEBUG, format, args...)
}
