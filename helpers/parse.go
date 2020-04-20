package helpers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"

	"github.com/golang/glog"
)

func ParseJSON(jsonStream *io.ReadCloser, outputStruct interface{}) error {
	body, err := ioutil.ReadAll(*jsonStream)

	if err != nil {
		glog.Errorln("ParseJSON ReadAll failed")
		glog.Errorln(err.Error())
		return err
	}

	if err = json.Unmarshal(body, outputStruct); err != nil {
		glog.Errorln("ParseJSON Unmarshal failed")
		glog.Errorln(err.Error())
		return err
	}

	return nil
}

func ParseIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		glog.Infof("error getting ip: %s", err.Error())
		return "0.0.0.0"
	}

	return ip
}
