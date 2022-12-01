package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var isSystemd bool

var logLevel int

func init() {
	// if val, ok := os.LookupEnv("INVOCATION_ID"); ok {
	// 	if len(val) > 0 {
	// 		isSystemd = true
	// 	}
	// }

	parentpid := os.Getppid()
	if parentpid == 1 {
		isSystemd = true
	}

	if val, ok := os.LookupEnv("AH_LOGLEVEL"); ok {
		if val == "debug" {
			logLevel = 1
		}
		if val == "trace" {
			logLevel = 2
		}
	}
}

func GetDebugLevel() int {
	if val, ok := os.LookupEnv("AH_LOGLEVEL"); ok {
		if val == "debug" {
			return 1
		}
	}
	return 0
}

func Info(msg ...interface{}) {
	write("[INFO] ", msg...)
}

func Infof(msg string, any ...interface{}) {
	writef("[INFO] ", msg, any...)
}

func Panic(msg ...interface{}) {
	write("[PANIC] ", msg...)
	panic(msg)
}

func Error(msg ...interface{}) {
	write("[ERROR] ", msg...)
}

func Errorf(msg string, any ...interface{}) {
	writef("[ERROR] ", msg, any...)
}

func Warning(msg ...interface{}) {
	write("[WARN] ", msg...)
}

func Debug(msg ...interface{}) {
	if logLevel > 0 {
		write("[DEBUG] ", msg...)
	}
}

func Trace(msg ...interface{}) {
	if logLevel > 1 {
		write("[TRACE] ", msg...)
	}
}

func write(logPrefix string, message ...interface{}) {
	msg := fmt.Sprintln(message...)
	print(logLevel, logPrefix, msg)
}

func writef(logPrefix string, message string, any ...interface{}) {
	msg := fmt.Sprintf(message, any...)
	print(logLevel, logPrefix, msg)
}

func print(level int, logPrefix string, msg string) {
	var prefix string
	if level > 0 {
		pc := make([]uintptr, 5)
		runtime.Callers(0, pc)

		f := runtime.FuncForPC(pc[4])
		location := f.Name()

		prefix = logPrefix + "F:" + location + ":"
	} else {
		prefix = ""
	}

	if !isSystemd {
		fmt.Print(time.Now().Local().Format("2006/01/02 15:04:05 "))
	}
	fmt.Print(prefix, msg)
}
