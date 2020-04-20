package mls

import (
	"fmt"
	"net/http"

	db "github.com/MarcusOuelletus/rets-server/database"
	"github.com/MarcusOuelletus/rets-server/helpers"
	"github.com/MarcusOuelletus/rets-server/templates"

	"github.com/golang/glog"
)

func MLS(r *http.Request) *templates.APIResponse {
	var response = &templates.APIResponse{
		Response: nil,
		Error:    "",
	}

	var result templates.Property
	var params map[string]interface{}

	helpers.ParseJSON(&r.Body, &params)

	err := fetchResults(params, &result)
	if err != nil {
		glog.Errorf("Search/Select: %s\n", err.Error())
		response.Error = "error fetching results"
		return response
	}

	response.Response = result

	return response
}

func fetchResults(params map[string]interface{}, result *templates.Property) error {
	var query = &db.Query{
		Conditions: map[string]interface{}{
			templates.Fields.MLS: params["mls"],
		},
		FieldsToReturn: templates.SearchReturn,
		Destination:    result,
		Page:           1,
	}

	for _, value := range []string{"res", "con", "com"} {

		if params["sold"] != nil && params["sold"].(bool) {
			query.Collection = fmt.Sprintf("%s-sold", value)
		} else {
			query.Collection = value
		}

		if err := db.SelectRow(query); err == nil {
			break
		}

	}

	return nil
}
