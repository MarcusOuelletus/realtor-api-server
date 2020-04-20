package logs

import (
	"flag"

	"github.com/golang/glog"
)

const logFileRelativePath = "logs"

func InitDevLogging() {
	flag.Set("stderrthreshold", "INFO")
	glog.Infoln("log level set to dev")
}
