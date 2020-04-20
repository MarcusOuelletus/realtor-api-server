package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"

	"github.com/MarcusOuelletus/rets-server/global"
	"github.com/MarcusOuelletus/rets-server/templates"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	Collection     string
	Conditions     map[string]interface{}
	FieldsToReturn *[]string
	Destination    interface{}
	Page           int
	OrderBy        string
	DistinctField  string
	Limit          int
	Hint           string
}

type InsertQuery struct {
	Collection string
	Data       map[string]interface{}
}

type UpsertQuery struct {
	Collection string
	Conditions map[string]interface{}
	Data       map[string]interface{}
}

type DistinctQuery struct {
	Collection     string
	Conditions     map[string]interface{}
	FieldsToReturn *[]string
	Destination    *[]interface{}
	Page           int
	OrderBy        string
	DistinctField  string
	Limit          int
}

type db struct {
	clientToDisconnect *mongo.Client
	collection         *mongo.Collection
}

var err error

func OpenDatabase() (*mongo.Database, error) {
	var clientOptions = &options.ClientOptions{
		Auth: &options.Credential{
			AuthSource: "AUTH_SOURCE",
			Username:   "USERNAME",
			Password:   "PASSWORD",
		},
	}

	client, err := mongo.Connect(
		context.Background(),
		clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:27017", global.DatabaseIP)),
	)

	if err != nil {
		glog.Errorln("error - creating knew mongo client")
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		glog.Errorln("error - pinging mongo connection")
		return nil, err
	}

	d := client.Database(global.DATABASE_NAME)

	return d, nil
}

func openSession(requestedCollection string) (*db, error) {
	var clientOptions = &options.ClientOptions{
		Auth: &options.Credential{
			AuthSource: "AUTH_SOURCE",
			Username:   "USERNAME",
			Password:   "PASSWORD",
		},
	}
	client, err := mongo.Connect(
		context.Background(),
		clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:27017", global.DatabaseIP)))
	if err != nil {
		glog.Errorln("error - creating knew mongo client")
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		glog.Errorln("error - pinging mongo connection")
	}

	collection := client.Database(global.DATABASE_NAME).Collection(requestedCollection)

	return &db{
		clientToDisconnect: client,
		collection:         collection,
	}, nil
}

func concatFieldsToReturn(fields *[]string) bson.M {
	var b = make(bson.M, len(*fields))

	for _, v := range *fields {
		b[v] = 1
	}

	return b
}

func concatConditions(conditions map[string]interface{}) bson.M {
	var b = make(bson.M, len(conditions))

	for k, v := range conditions {
		if v == "" {
			continue
		}
		b[k] = v
	}

	return b
}

func Select(query *Query) error {
	db, err := openSession(query.Collection)
	if err != nil {
		return errors.New("Error Connecting to Database")
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	var limit = global.QUERY_LIMIT

	if query.Limit != 0 {
		limit = query.Limit
	}

	var queryOptions = &options.FindOptions{
		Projection: concatFieldsToReturn(query.FieldsToReturn),
	}

	queryOptions = queryOptions.SetLimit(int64(limit))
	queryOptions = queryOptions.SetSkip(int64(query.Page * global.QUERY_LIMIT))

	if query.Hint != "" {
		queryOptions = queryOptions.SetHint(bson.M{query.Hint: 1})
	}

	if query.OrderBy != "" {
		queryOptions.SetSort(bson.M{query.OrderBy: -1})
	}

	var conditionsFilter = primitive.M{}

	if query.Conditions != nil {
		query.Conditions[templates.Fields.MarcPix] = HasValue()

		if query.Collection == "com" {
			delete(query.Conditions, "Xd")
		}

		conditionsFilter = concatConditions(query.Conditions)
	}

	cursor, err := db.collection.Find(context.Background(), conditionsFilter, queryOptions)
	if err != nil {
		glog.Errorln("error performing Select, Find() failed")
		return err
	}

	if err := cursor.All(context.Background(), query.Destination); err != nil {
		glog.Errorln("error performing Select, All() failed")
		return err
	}

	return nil
}

func SelectDistinctWithoutCondition(query *DistinctQuery) error {
	db, err := openSession(query.Collection)
	if err != nil {
		return errors.New("Error Connecting to Database")
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	*query.Destination, err = db.collection.Distinct(context.Background(), query.DistinctField, bson.M{query.DistinctField: bson.M{"$ne": nil}})

	return nil
}

func SelectDistinct(query *DistinctQuery) error {
	db, err := openSession(query.Collection)
	if err != nil {
		return errors.New("Error Connecting to Database")
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	*query.Destination, err = db.collection.Distinct(context.Background(), query.DistinctField, bson.M{})

	return nil
}

func SelectRow(query *Query) error {
	db, err := openSession(query.Collection)
	if err != nil {
		glog.Errorln("Error Connecting to Database")
		return err
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	var cursor *mongo.SingleResult

	if query.FieldsToReturn != nil {
		findOneOptions := &options.FindOneOptions{
			Projection: concatFieldsToReturn(query.FieldsToReturn),
		}
		cursor = db.collection.FindOne(context.Background(), concatConditions(query.Conditions), findOneOptions)
	} else {
		cursor = db.collection.FindOne(context.Background(), concatConditions(query.Conditions))
	}

	if err = cursor.Decode(query.Destination); err != nil {
		glog.Errorf("error - SelectRow/Decode - %s\n", err.Error())
		return err
	}

	return nil
}

func Upsert(query *UpsertQuery) error {
	db, err := openSession(query.Collection)
	if err != nil {
		glog.Errorln("Error Connecting to Database")
		return err
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	var upsertOptions = &options.UpdateOptions{}
	upsertOptions = upsertOptions.SetUpsert(true)

	_, err = db.collection.UpdateOne(
		context.Background(),
		query.Conditions,
		map[string]interface{}{
			"$set": query.Data,
		},
		upsertOptions,
	)

	return err
}

func Update(query *UpsertQuery) error {
	db, err := openSession(query.Collection)
	if err != nil {
		glog.Errorln("Error Connecting to Database")
		return err
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	_, err = db.collection.UpdateOne(
		context.Background(),
		query.Conditions,
		map[string]interface{}{
			"$set": query.Data,
		},
		nil,
	)

	return err
}

func Insert(query *InsertQuery) error {
	db, err := openSession(query.Collection)
	if err != nil {
		glog.Errorln("Error Connecting to Database")
		return err
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	_, err = db.collection.InsertOne(context.Background(), query.Data)
	if err != nil {
		glog.Errorln("Error Inserting Into Database")
		return err
	}

	return nil
}

func Remove(query *InsertQuery) error {
	db, err := openSession(query.Collection)
	if err != nil {
		glog.Errorln("Error Connecting to Database")
		return err
	}

	defer db.clientToDisconnect.Disconnect(context.Background())

	_, err = db.collection.DeleteMany(context.Background(), query.Data)

	return err
}
