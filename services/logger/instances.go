package logger

import (
	"os"
)

func Default() Logger {
	if stdLogger != nil {
		return stdLogger
	}
	Config()

	return stdLogger
}

// var HttpLogger Logger
var stdLogger Logger

// log file
var httpFile *os.File

func Config() {
	buildLoggerInstances()
}

func buildLoggerInstances() {
	stdLogger = NewLogDirector(&StdLoggerBuilder{}).BuildLogger()
}
