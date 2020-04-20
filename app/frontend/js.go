package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/MarcusOuelletus/rets-server/app/access/tokens"
	"github.com/MarcusOuelletus/rets-server/global"
	"github.com/MarcusOuelletus/rets-server/helpers/brokers"

	"github.com/golang/glog"
)

type TemplateData struct {
	Tables       []string `json:"tables"`
	Token        string   `json:"token"`
	ColorPrimary string   `json:"colorPrimary"`
	Logo         string   `json:"logo"`
	All          bool     `json:"all"`
	Sold         bool     `json:"sold"`
	Areas        []string `json:"areas"`
}

type tokenResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

func JS(w http.ResponseWriter, r *http.Request) {
	var err error

	templateData := &TemplateData{}

	if err := r.ParseForm(); err != nil {
		glog.Errorf("error parsing post form values: %s\n", err.Error())
		return
	}

	jsonBytes := []byte(r.PostFormValue("data"))

	if err := json.Unmarshal(jsonBytes, &templateData); err != nil {
		glog.Errorf("error unmarshalling post params: %s\n", err.Error())
		return
	}

	brokerToken := templateData.Token

	broker, err := brokers.Get(brokerToken)

	if err != nil {
		return
	}

	templateData.Tables = broker.Tables

	// set token to encrypted client token
	templateData.Token, err = tokens.CreateClientToken(brokerToken, r.RemoteAddr)

	if err != nil {
		return
	}

	htmlTemplatePath := fmt.Sprintf("%s/app/frontend/index.html", global.WorkingDirectory)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		glog.Errorf("error creating frontend template: %s\n", err.Error())
		return
	}

	if err = tmpl.ExecuteTemplate(w, "index.html", templateData); err != nil {
		glog.Errorf("error executing frontend template: %s\n", err.Error())
		return
	}
}
