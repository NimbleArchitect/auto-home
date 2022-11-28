package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var isSystemd bool

type logger struct {
	// 0 = errors only, 1 = errors + debug, 2 = trace + debug + errors
	level *int
}

func init() {
	// if val, ok := os.LookupEnv("INVOCATION_ID"); ok {
	// 	if len(val) > 0 {
	// 		isSystemd = true
	// 	}
	// }
}

// TODO: need to automatically pickup the function name the functions are called from

func New(level *int) *logger {
	return &logger{level: level}
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
	l.write("[INFO] ", msg...)
}

func (l *logger) Infof(msg string, any ...interface{}) {
	l.writef("[INFO] ", msg, any...)
}

func (l *logger) Panic(msg ...interface{}) {
	l.write("[PANIC] ", msg...)
	panic(msg)
}

func (l *logger) Error(msg ...interface{}) {
	l.write("[ERROR] ", msg...)
}

func (l *logger) Errorf(msg string, any ...interface{}) {
	l.writef("[ERROR] ", msg, any...)
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
	msg := fmt.Sprintln(message...)
	print(*l.level, logPrefix, msg)
}

func (l *logger) writef(logPrefix string, message string, any ...interface{}) {
	msg := fmt.Sprintf(message, any...)
	print(*l.level, logPrefix, msg)
}

func print(level int, logPrefix string, msg string) {
	var prefix string
	if level > 0 {
		pc := make([]uintptr, 2)
		runtime.Callers(3, pc)
		f := runtime.FuncForPC(pc[0])
		location := f.Name()
		if location == "server/logger.(*logger).Debug" {
			f := runtime.FuncForPC(pc[1])
			location = f.Name()
		}

		path := strings.Split(location, "/")
		location = path[len(path)-1]

		prefix = logPrefix + "F:" + location + ":"
	} else {
		prefix = ""
	}

	if !isSystemd {
		fmt.Print(time.Now().Local().Format("2006/01/02 15:04:05 "))
	}
	fmt.Print(prefix, msg)
}
