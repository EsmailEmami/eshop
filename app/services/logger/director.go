package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type logDirector struct {
	logBuilder LogBuilder
}

func (d *logDirector) BuildLogger() Logger {
	d.logBuilder.AddWriters()
	ll := newLogger(d.logBuilder.GetWriters()...)

	return ll
}

func NewLogDirector(builder LogBuilder) *logDirector {
	return &logDirector{
		logBuilder: builder,
	}
}

func newLogger(writers ...io.Writer) *logger {
	var logWriter io.Writer

	if len(writers) == 0 {
		logWriter = os.Stdout
	} else {
		mw := io.MultiWriter(writers...)
		logWriter = mw
	}
	l := zerolog.New(logWriter).With().Timestamp().Logger()
	ll := logger{}
	ll.LoggerInstance = &l
	ll.Writers = append(ll.Writers, writers...)
	ll.Fields = make(map[string]string)

	return &ll
}
