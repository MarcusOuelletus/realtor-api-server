package database

import (
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func LessOrEqual(input string) bson.M {
	var value float64
	var err error
	if input == "" {
		value = 1000000000000
	} else {
		value, err = strconv.ParseFloat(input, 32)
		if err != nil {
			log.Printf("error - Search/lessOrEqual - %s", err.Error())
			return bson.M{}
		}
	}
	return bson.M{"$lte": value}
}

func GreaterOrEqual(input string) bson.M {
	var value float64
	var err error
	if input == "" {
		value = 0
	} else {
		value, err = strconv.ParseFloat(input, 32)
		if err != nil {
			log.Printf("error - Search/greaterAndLess - %s", err.Error())
			return bson.M{}
		}
	}
	return bson.M{"$gte": value}
}

func GreaterAndLess(inputMin, inputMax string) bson.M {
	var max, min float64
	var err error
	if inputMax == "" {
		max = 1000000000000.00
	} else {
		max, err = strconv.ParseFloat(inputMax, 32)
		if err != nil {
			log.Printf("error - Search/greaterAndLess - %s", err.Error())
			return bson.M{}
		}
	}
	if inputMin == "" {
		min = 0.00
	} else {
		min, err = strconv.ParseFloat(inputMin, 32)
		if err != nil {
			log.Printf("error - Search/greaterAndLess - %s", err.Error())
			return bson.M{}
		}
	}
	return bson.M{"$gte": min, "$lte": max}
}

func ClosestMatch(input string) bson.M {
	var wildcardString = fmt.Sprintf(`%s`, input)
	return bson.M{"$regex": wildcardString}
}

func HasValue() bson.M {
	return bson.M{"$ne": nil}
}

func HasPrefix(input string) bson.M {
	var wildcardString = fmt.Sprintf(`^%s`, input)
	return bson.M{"$regex": wildcardString}
}

func HasNotPrefix(input string) bson.M {
	var wildcardString = fmt.Sprintf(`^(?!%s)`, input)
	return bson.M{"$regex": wildcardString}
}

func NotEqualTo(input string) bson.M {
	return bson.M{"$ne": input}
}
