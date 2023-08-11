package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Warning(msg string)
	Error(msg string)
	Errorf(format string, args ...interface{})
	Fatal(msg string)
	Panic(msg string)
	WithField(key string, value string) Logger
	WithFields(fields map[string]string) Logger
}

type logger struct {
	LoggerInstance *zerolog.Logger
	Fields         map[string]string
	Writers        []io.Writer
}

func (entry *logger) bindFields(loggerEvent *zerolog.Event) *zerolog.Event {
	for k, v := range entry.Fields {
		loggerEvent.Str(k, v)
	}

	return loggerEvent
}

func (entry *logger) Info(msg string) {
	entry.bindFields(entry.LoggerInstance.Info()).Msg(msg)
}

func (entry *logger) Infof(format string, args ...interface{}) {
	entry.bindFields(entry.LoggerInstance.Info()).Msgf(format, args...)
}

func (entry *logger) Warning(msg string) {
	entry.bindFields(entry.LoggerInstance.Warn()).Msg(msg)
}

func (entry *logger) Error(msg string) {
	entry.bindFields(entry.LoggerInstance.Error()).Msg(msg)
}

func (entry *logger) Errorf(format string, args ...interface{}) {
	entry.bindFields(entry.LoggerInstance.Error()).Msgf(format, args...)
}

func (entry *logger) Fatal(msg string) {
	entry.bindFields(entry.LoggerInstance.Fatal()).Msg(msg)
}

func (entry *logger) Panic(msg string) {
	entry.bindFields(entry.LoggerInstance.Panic()).Msg(msg)
}

func (entry *logger) WithField(key string, value string) Logger {
	newEntry := newLogger(entry.Writers...)
	for k, v := range entry.Fields {
		newEntry.Fields[k] = v
	}
	newEntry.Fields[key] = value
	return newEntry
}

func (entry *logger) WithFields(fields map[string]string) Logger {
	newEntry := newLogger(entry.Writers...)
	for k, v := range entry.Fields {
		newEntry.Fields[k] = v
	}

	for k, v := range fields {
		newEntry.Fields[k] = v
	}
	return newEntry
}
