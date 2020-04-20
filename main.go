package main

import (
	"flag"
	"os"

	"github.com/MarcusOuelletus/rets-server/app/areas"
	"github.com/MarcusOuelletus/rets-server/logs"
	"github.com/MarcusOuelletus/rets-server/server"

	"github.com/golang/glog"
)

func main() {
	parseFlags()

	defer glog.Flush()

	areas.BuildCache()

	server.Start()
}

func parseFlags() {
	isDev := flag.Bool("log", false, "boolean flag: default is no logs")
	flag.Parse()

	os.Setenv("MLS_BASE_URL", "http://localhost:9090")

	if *isDev {
		logs.InitDevLogging()
	}
}
