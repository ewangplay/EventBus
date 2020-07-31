package log

import (
	"fmt"
	"path/filepath"

	"github.com/arxanfintech/go-logger"
	"github.com/ewangplay/eventbus/config"
)

type Logger struct {
	service_name  string
	cfgLevel      config.LogLevel
	logger        *logger.Logger
	rotate_writer *RotateWriter
}

func New(opts *config.EB_Options) (log *Logger, err error) {
	service := opts.ServiceName
	filename := fmt.Sprintf("%s.log", service)
	log_file := filepath.Join(opts.LogPath, filename)

	rotate_writer, err := NewRotateWriter(log_file,
		opts.LogMaxSize,
		opts.LogRotateDaily)
	if err != nil {
		return nil, err
	}

	color := 0
	logger, err := logger.New(service, color, rotate_writer)
	if err != nil {
		return nil, err
	}

	log = &Logger{
		service_name:  service,
		cfgLevel:      config.ParseLogLevel(opts.LogLevel),
		logger:        logger,
		rotate_writer: rotate_writer,
	}

	return log, nil
}

func (this *Logger) Close() (err error) {
	return this.rotate_writer.Close()
}

func (this *Logger) log(level config.LogLevel, format string, args ...interface{}) (err error) {
	if level < this.cfgLevel {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	fmt.Println(msg)

	s := fmt.Sprintf("%s: %s", level, msg)

	this.Output(3, s)

	return nil
}

func (this *Logger) Output(maxdepth int, s string) error {
	maxdepth += 2
	this.logger.Log("", maxdepth, s)
	return nil
}

func (this *Logger) Fatal(format string, args ...interface{}) (err error) {
	return this.log(config.FATAL, format, args...)
}

func (this *Logger) Error(format string, args ...interface{}) (err error) {
	return this.log(config.ERROR, format, args...)
}

func (this *Logger) Warn(format string, args ...interface{}) (err error) {
	return this.log(config.WARN, format, args...)
}

func (this *Logger) Info(format string, args ...interface{}) (err error) {
	return this.log(config.INFO, format, args...)
}

func (this *Logger) Debug(format string, args ...interface{}) (err error) {
	return this.log(config.DEBUG, format, args...)
}
