package search

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MarcusOuelletus/rets-server/app/access"
	"github.com/MarcusOuelletus/rets-server/app/access/tokens"
	db "github.com/MarcusOuelletus/rets-server/database"
	"github.com/MarcusOuelletus/rets-server/global"
	"github.com/MarcusOuelletus/rets-server/helpers"
	"github.com/MarcusOuelletus/rets-server/helpers/brokers"
	"github.com/MarcusOuelletus/rets-server/templates"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var err error

type search struct {
	broker   *brokers.Broker
	response *templates.APIResponse
}

// Handler - Called upon a user search
func Handler(r *http.Request) *templates.APIResponse {
	var s = search{
		response: &templates.APIResponse{
			Token:    "",
			Response: nil,
			Error:    "",
			Code:     0,
		},
	}

	tokenData, err := tokens.Parse(r)

	if err != nil {
		s.response.Error = templates.InvalidToken
		return s.response
	}

	s.broker, err = brokers.Get(tokenData.BrokerToken)

	if err != nil {
		s.response.Error = templates.InvalidToken
		return s.response
	}

	requestData, err := s.parseRequestJSON(&r.Body)

	if err != nil {
		s.response.Error = "error parsing json from request"
		return s.response
	}

	properties, err := s.fetchResults(requestData)

	if err != nil {
		s.response.Error = "error fetching results"
		return s.response
	}

	s.response.Response = properties

	return s.response
}

func (s *search) parseRequestJSON(requestBody *io.ReadCloser) (*templates.SearchRequest, error) {
	var requestData = new(templates.SearchRequest)

	if err := helpers.ParseJSON(requestBody, requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func (s *search) fetchResults(requestData *templates.SearchRequest) ([]templates.Property, error) {
	var properties []templates.Property

	conditions := s.buildConditionsMap(requestData)

	var query = &db.Query{
		Collection:     requestData.Type,
		Conditions:     conditions,
		FieldsToReturn: templates.SearchReturn,
		Page:           requestData.Page,
		OrderBy:        templates.Fields.Timestamp,
		Hint:           templates.Fields.Brokerage,
	}

	if requestData.Sold {
		if err := s.restrictSoldAccess(requestData.UserID, s.broker); err != nil {
			return nil, err
		}

		s.modifyQueryForSoldProperties(query)
	}

	if requestData.Area != "" {
		s.modifyQueryForSpecificLocation(query, requestData.Area)
	}

	// Always have specific brokerage be on top
	results, err := s.querySpecificBrokerage(query, s.broker.Brokerage, requestData.Area)

	if err != nil {
		s.response.Code = 100
		glog.Errorln(err.Error())
		return nil, err
	}

	properties = append(properties, results...)

	// If brokerage doesn't have enough properties to fill need, query other brokerages
	if !requestData.All || len(properties) >= global.QUERY_LIMIT {
		return properties, nil
	}

	query.Limit = global.QUERY_LIMIT - len(properties)

	results, err = s.queryOtherThanBrokerage(query, s.broker.Brokerage, requestData.Area)

	if err != nil {
		s.response.Code = 100
		glog.Errorln(err.Error())
		return nil, err
	}

	properties = append(properties, results...)

	return properties, nil
}

func (s *search) restrictSoldAccess(userID string, broker *brokers.Broker) error {
	if !broker.SoldAccess {
		s.response.Code = 801
		return fmt.Errorf("broker %s hasn't purchased sold access", broker.Name)
	}

	if userID == "" || !access.UserIDExists(userID) {
		s.response.Code = 800
		return errors.New("user doesn't exist")
	}

	return nil
}

func (s *search) modifyQueryForSoldProperties(query *db.Query) {
	query.Collection = fmt.Sprintf("%s-sold", query.Collection)

	query.Conditions["$nor"] = bson.A{
		bson.M{"Sp_dol": nil},
		bson.M{"Sp_dol": ""},
	}

	query.Conditions[templates.Fields.Status] = "U"
}

func (s *search) modifyQueryForSpecificLocation(query *db.Query, location string) error {
	if v, ok := query.Conditions["$or"]; ok {
		glog.Errorln("can't have multiple $or commands")
		glog.Infof("$or is already %+v\n", v)
		return errors.New("can't have multiple $or commands")
	}

	locationRegex := fmt.Sprintf("^%s", location)

	query.Conditions["$or"] = bson.A{
		bson.M{templates.Fields.Area: primitive.Regex{Pattern: locationRegex, Options: ""}},
		bson.M{templates.Fields.Community: primitive.Regex{Pattern: locationRegex, Options: ""}},
		bson.M{templates.Fields.Municipality: primitive.Regex{Pattern: locationRegex, Options: ""}},
	}

	return nil
}

func (s *search) querySpecificBrokerage(query *db.Query, brokerage, location string) ([]templates.Property, error) {
	query.Conditions[templates.Fields.Brokerage] = brokerage

	return s.query(query, location)
}

func (s *search) queryOtherThanBrokerage(query *db.Query, brokerage, location string) ([]templates.Property, error) {
	query.Conditions[templates.Fields.Brokerage] = db.NotEqualTo(brokerage)

	return s.query(query, location)
}

func (s *search) query(query *db.Query, location string) ([]templates.Property, error) {
	var results []templates.Property

	query.Destination = &results

	if err := db.Select(query); err != nil {
		return nil, err
	}

	return results, nil
}
