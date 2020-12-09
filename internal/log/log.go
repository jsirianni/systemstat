package log

import (
	"flag"

	"github.com/golang/glog"
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Parse()
}

func Info(m string) {
	glog.Infof(m)
	glog.Flush()
}

func Error(err error) {
	glog.Error(err)
	glog.Flush()
}
