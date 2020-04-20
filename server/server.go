package server

import (
	"net/http"

	"github.com/golang/glog"
)

func Start() {
	http.Handle("/", generateRouter())

	if err := http.ListenAndServe(":9090", nil); err != nil {
		glog.Infof("error starting server: %s\n", err.Error())
	}
}
