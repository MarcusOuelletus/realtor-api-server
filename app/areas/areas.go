package areas

import (
	"io"
	"net/http"
	"time"

	"github.com/MarcusOuelletus/rets-server/database"
	"github.com/MarcusOuelletus/rets-server/global"
	"github.com/MarcusOuelletus/rets-server/helpers"
	"github.com/MarcusOuelletus/rets-server/templates"
)

type areasInstance struct {
	query                    *database.Query
	distinctQuery            *database.DistinctQuery
	requestData              *areaParams
	results                  []interface{}
	currentFieldBeingFetched string
}

type areaParams struct {
	PropertyType string `json:"type"`
	Value        string `json:"value"`
	Page         int    `json:"page"`
}

type areasCacheObject struct {
	timestamp      time.Time
	cachedAllAreas []interface{}
}

var areasCache = &areasCacheObject{
	timestamp:      time.Now(),
	cachedAllAreas: nil,
}

var err error

func Areas(r *http.Request) *templates.APIResponse {
	var response = &templates.APIResponse{
		Response: nil,
		Error:    "",
	}

	var a areasInstance

	a.requestData, err = a.parseRequestJSON(&r.Body)

	if err != nil {
		response.Error = "error found in request JSON"
		return response
	}

	results := QueryCache(a.requestData.Value, a.requestData.Page*global.AREA_LIMIT)

	response.Response = results

	return response
}

func (a *areasInstance) parseRequestJSON(requestBody *io.ReadCloser) (*areaParams, error) {
	var requestData = new(areaParams)

	if err := helpers.ParseJSON(requestBody, requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}
