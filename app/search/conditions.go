package search

import (
	"time"

	"github.com/MarcusOuelletus/rets-server/database"

	"github.com/MarcusOuelletus/rets-server/templates"
	"go.mongodb.org/mongo-driver/bson"
)

type conditions struct {
	Map map[string]interface{}
}

func (s *search) buildConditionsMap(requestData *templates.SearchRequest) map[string]interface{} {
	var c = &conditions{
		Map: make(map[string]interface{}),
	}

	c.Map[templates.Fields.Status] = "A"
	c.Map[templates.Fields.DisplayAddress] = "Y"

	if !requestData.Sold && requestData.Type != "com" {
		c.Map[templates.Fields.Expiry] = bson.M{"$gt": time.Now().Format("2006-01-02 15:04:05")}
	}

	if requestData.MinPrice != "" || requestData.MaxPrice != "" {
		c.Map[templates.Fields.Price] = database.GreaterAndLess(requestData.MinPrice, requestData.MaxPrice)
	}

	if requestData.MLS != "" {
		c.Map[templates.Fields.MLS] = requestData.MLS
	}

	if requestData.Address != "" {
		c.Map[templates.Fields.Address] = database.ClosestMatch(requestData.Address)
	}

	if requestData.Type != "com" {
		c.addResAndConConditions(requestData)
	}

	return c.Map
}

func (c *conditions) addResAndConConditions(requestData *templates.SearchRequest) {
	if requestData.Beds != "" {
		c.Map[templates.Fields.Beds] = database.GreaterOrEqual(requestData.Beds)
	}

	if requestData.Baths != "" {
		c.Map[templates.Fields.Baths] = database.GreaterOrEqual(requestData.Baths)
	}

	if requestData.BuildingType != "" {
		c.Map[templates.Fields.BuildingType] = requestData.BuildingType
	}

	if requestData.SaleType != "" {
		c.Map[templates.Fields.SaleType] = requestData.SaleType
	}

	if requestData.GarageType != "" {
		c.Map[templates.Fields.GarageType] = requestData.GarageType
	}

	if requestData.YearBuilt != "" {
		c.Map[templates.Fields.YearBuilt] = requestData.YearBuilt
	}
}
