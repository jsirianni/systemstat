package log

import (
	"os"
	"flag"

	"github.com/golang/glog"
)

func init() {
	flag.Set("logtostderr", "true")
	// stderrthreshold=[WARNING|ERROR|INFO|FATAL]
	flag.Set("stderrthreshold", "WARNING")
//	flag.Parse()
}

func Info(m ...string) {
	glog.Infoln(m)
	glog.Flush()
}

func Error(err error) {
	glog.Error(err)
	glog.Flush()
}

func Fatal(err error, exitCode int) {
	Error(err)
	os.Exit(exitCode)
}

func Trace(err error) {
	Error(err)
}
