package log

import (
	"os"
	"flag"

	"github.com/golang/glog"
)

const envLogLevl = "SYSTEMSTAT_LOG_LEVEL"

const (
	errorLVL = "ERROR"
	infoLVL  = "INFO"
	debugLVL = "DEBUG"
	traceLVL = "TRACE"

	errorLevel = 1
	infoLevel  = 2
	debugLevel = 3
	traceLevel = 4
)

var logLevel string
var level int

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")

	SetLogLevel(os.Getenv(envLogLevl))
}

func SetLogLevel(l string) {
	switch l {
	case errorLVL:
		logLevel = l
		level    = errorLevel
	case infoLVL:
		logLevel = l
		level    = infoLevel
	case debugLVL:
		logLevel = l
		level    = debugLevel
	case traceLVL:
		logLevel = l
		level    = traceLevel
	}
}

func Info(m ...string) {
	if level >= infoLevel {
		glog.Infoln(m)
		glog.Flush()
	}
}

func Error(err error) {
	if level >= errorLevel {
		glog.Error(err)
		glog.Flush()
	}
}

// Error, except exit with the given status code
func Fatal(err error, exitCode int) {
	glog.Error(err)
	glog.Flush()
	os.Exit(exitCode)
}

func Debug(err error) {
	if level >= debugLevel {
		glog.Error(err)
		glog.Flush()
	}
}

func Trace(m ...string) {
	if level >= traceLevel {
		glog.Infoln(m)
		glog.Flush()
	}
}
