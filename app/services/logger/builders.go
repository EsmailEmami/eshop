package logger

import (
	"io"
	"os"
)

type LogBuilder interface {
	AddWriters()
	GetWriters() []io.Writer
	GetLogger() logger
}

type HttpLogBuilder struct {
	logger
}

func (l *HttpLogBuilder) AddWriters() {
	l.Writers = append(l.Writers, os.Stdout, httpFile)
}

func (l *HttpLogBuilder) GetWriters() []io.Writer {
	return l.Writers
}

func (l *HttpLogBuilder) GetLogger() logger {
	return l.logger
}

type StdLoggerBuilder struct {
	logger
}

func (l *StdLoggerBuilder) AddWriters() {
	l.Writers = append(l.Writers, os.Stdout)
}

func (l *StdLoggerBuilder) GetWriters() []io.Writer {
	return l.Writers
}

func (l *StdLoggerBuilder) GetLogger() logger {
	return l.logger
}
