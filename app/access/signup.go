package access

import (
	"fmt"
	"net/http"

	"github.com/MarcusOuelletus/rets-server/app/access/tokens"
	"github.com/MarcusOuelletus/rets-server/database"
	"github.com/MarcusOuelletus/rets-server/email"
	"github.com/MarcusOuelletus/rets-server/helpers"
	"github.com/MarcusOuelletus/rets-server/helpers/brokers"
	"github.com/MarcusOuelletus/rets-server/templates"

	"github.com/golang/glog"
)

type userData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type signupInstace struct {
	BrokerToken string
	Brokerage   string
	Login       bool
	User        *userData
	Response    *templates.APIResponse
}

func Signup(r *http.Request) *templates.APIResponse {
	var signup = &signupInstace{
		Login:       false,
		User:        &userData{},
		BrokerToken: "",
		Response: &templates.APIResponse{
			Response: nil,
			Error:    "",
			Code:     0,
		},
	}

	helpers.ParseJSON(&r.Body, signup.User)

	var tokenData, err = tokens.Parse(r)

	if err != nil {
		signup.Response.Code = 1
		return signup.Response
	}

	signup.BrokerToken = tokenData.BrokerToken

	if signup.User.FirstName == "" || signup.User.LastName == "" {
		signup.Login = true
	}

	if signup.Login {
		signup.Response.Response = signup.handleLogin()
	} else {
		signup.Response.Response = signup.handleSignup()
	}

	return signup.Response
}

func (s *signupInstace) handleLogin() string {
	id := s.getIDFromEmail()
	if id == "" {
		s.Response.Code = 1
	}
	return id
}

func (s *signupInstace) handleSignup() string {
	broker, err := brokers.Get(s.BrokerToken)

	if err != nil {
		s.Response.Code = 1
		return ""
	}

	id, err := s.createUser(broker.Token)

	if err != nil {
		glog.Errorf("error inserting user: %s\n", err.Error())
		s.Response.Code = 1
		return ""
	}

	email.SendEmail(&email.EmailObject{
		Recipient: broker.Email,
		BodyText:  fmt.Sprintf("New Lead: %s %s - %s", s.User.FirstName, s.User.LastName, s.User.Email),
	})

	return id
}

func (s *signupInstace) createUser(brokerToken string) (string, error) {
	glog.Infof("creating user %s\n", s.User.Email)
	var id = helpers.GenerateToken(30)

	return id, database.Upsert(&database.UpsertQuery{
		Collection: "users",
		Conditions: map[string]interface{}{
			"email": s.User.Email,
		},
		Data: map[string]interface{}{
			"first_name": s.User.FirstName,
			"last_name":  s.User.LastName,
			"email":      s.User.Email,
			"id":         id,
			"broker":     brokerToken,
		},
	})
}

func (s *signupInstace) getIDFromEmail() string {
	glog.Infof("checking if user %s exists\n", s.User.Email)

	var dbResult map[string]interface{}

	database.SelectRow(&database.Query{
		Collection: "users",
		Conditions: map[string]interface{}{
			"email": s.User.Email,
		},
		Destination: &dbResult,
	})

	if dbResult != nil {
		return dbResult["id"].(string)
	}

	return ""
}

func UserIDExists(id string) bool {
	glog.Infof("checking if user %s exists\n", id)

	var dbResult map[string]interface{}

	database.SelectRow(&database.Query{
		Collection: "users",
		Conditions: map[string]interface{}{
			"id": id,
		},
		Destination: &dbResult,
	})

	return dbResult != nil
}
