package brokers

import (
	"errors"

	"github.com/MarcusOuelletus/rets-server/database"
	"github.com/golang/glog"
)

type Broker struct {
	Brokerage  string   `bson:"Brokerage"`
	Token      string   `bson:"Token"`
	Email      string   `bson:"Email"`
	Tables     []string `bson:"Tables"`
	Name       string   `bson:"Name"`
	SoldAccess bool     `bson:"Sold"`
}

func Get(brokerToken string) (*Broker, error) {
	return getBrokerFromBrokerToken(brokerToken)
}

func getBrokerFromBrokerToken(brokerToken string) (*Broker, error) {
	var broker = &Broker{}

	var query = &database.Query{
		Collection: "accounts",
		Conditions: map[string]interface{}{
			"Token": brokerToken,
		},
		Destination: broker,
	}

	if err := database.SelectRow(query); err != nil {
		glog.Errorln("error setting broker using token")
		return nil, errors.New("error setting broker using token")
	}

	return broker, nil
}
