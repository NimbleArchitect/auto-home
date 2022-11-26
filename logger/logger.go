package logger

import (
	"fmt"
	"os"
)

type logger struct {
	location string
	// 0 = errors only, 1 = errors + debug, 2 = trace + debug + errors
	level *int
}

func New(location string, level *int) *logger {
	return &logger{location: location, level: level}
}

func GetDebugLevel() int {
	if val, ok := os.LookupEnv("AH_LOGLEVEL"); ok {
		if val == "debug" {
			return 1
		}
	}
	return 0
}

func (l *logger) Info(msg ...interface{}) {
	l.write("", msg...)
}

func (l *logger) Infof(msg string, any ...interface{}) {
	l.writef("", msg, any...)
}

func (l *logger) Panic(msg ...interface{}) {
	l.write("", msg...)
	panic(msg)
}

func (l *logger) Error(msg ...interface{}) {
	l.write("[ERROR] ", msg...)
}

func (l *logger) Warning(msg ...interface{}) {
	l.write("[WARN] ", msg...)
}

func (l *logger) Debug(msg ...interface{}) {
	if *l.level > 0 {
		l.write("[DEBUG] ", msg...)
	}
}

func (l *logger) Trace(msg ...interface{}) {
	if *l.level > 1 {
		l.write("[TRACE] ", msg...)
	}
}

func (l *logger) write(logPrefix string, message ...interface{}) {
	var msg string
	var prefix string

	msg = fmt.Sprintln(message...)

	if *l.level > 0 {
		prefix = logPrefix + "F:" + l.location + ":"
	} else {
		prefix = ""
	}

	fmt.Print(prefix, msg)
}

func (l *logger) writef(logPrefix string, msg string, any ...interface{}) {
	var prefix string

	msg = fmt.Sprintf(msg, any...)

	if *l.level > 0 {
		prefix = logPrefix + "F:" + l.location + ":"
	} else {
		prefix = ""
	}

	fmt.Print(prefix, msg)
}
