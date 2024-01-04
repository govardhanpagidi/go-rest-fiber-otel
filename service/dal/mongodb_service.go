package dal

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"

	"github.com/PeerIslands/aci-fx-go/model/dto/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	dbName         = "fx_data"
	collectionName = "forex_data"
)

type MongoDbService[T any] struct {
}

var database *mongo.Database

func (db *MongoDbService[T]) Init(credentials ...string) {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentServerSelection, options.LogLevelDebug)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(credentials[0]).SetLoggerOptions(loggerOptions))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database(dbName)
	return
}

func (db *MongoDbService[T]) GetDatabase() *mongo.Database {
	return database
}

func getContext() context.Context {
	return context.Background()
}

func (db *MongoDbService[T]) GetOne(filter any) (T, error) {
	var fxRequest = filter.(request.FxDataRequest)
	filterBson := bson.D{
		{"tenantId", fxRequest.TenantId},
		{"bankId", fxRequest.BankId},
		{"baseCurrency", fxRequest.BaseCurrency},
		{"targetCurrency", fxRequest.TargetCurrency},
		{"tier", fxRequest.Tier},
	}

	result := database.Collection(collectionName).FindOne(getContext(), filterBson)
	var data T
	err := result.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (db *MongoDbService[T]) GetOneById(id int) (T, error) {

	var data T

	return data, nil
}

func (db *MongoDbService[T]) Get(filter any) ([]T, error) {
	cursor, _ := database.Collection(collectionName).Find(getContext(), filter)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.Background())
	var data []T
	for cursor.Next(context.Background()) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			fmt.Println("Error decoding document:", err)
		}
		data = append(data, result)
	}
	return data, nil
}

func (db *MongoDbService[T]) CreateOne(document T) (T, error) {
	_, err := database.Collection(collectionName).InsertOne(getContext(), document)

	if err != nil {
		return document, err
	}
	return document, nil
}

func (db *MongoDbService[T]) UpdateOne(document any, filter any) (any, error) {
	var fxRequest = filter.(request.FxDataRequest)
	filterBson := bson.D{
		{"tenantId", fxRequest.TenantId},
		{"bankId", fxRequest.BankId},
		{"baseCurrency", fxRequest.BaseCurrency},
		{"targetCurrency", fxRequest.TargetCurrency},
		{"tier", fxRequest.Tier},
	}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	updateBson := bson.D{
		{"$inc", bson.D{
			{"buyRate", 0.01},
		}},
	}

	result := database.Collection(collectionName).FindOneAndUpdate(getContext(), filterBson, updateBson, option)

	var data T
	err := result.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (db *MongoDbService[T]) UpdateOneById(id any) (any, error) {
	var docVersion = id.(int)
	filterBson := bson.D{
		{"docVersion", docVersion},
	}
	updateBson := bson.D{
		{"$inc", bson.D{
			{"buyRate", 0.01},
		}},
	}

	result, err := database.Collection(collectionName).UpdateOne(getContext(), filterBson, updateBson)

	if err != nil || result.ModifiedCount == 0 {
		return false, err
	}
	return result.ModifiedCount, nil
}

func (db *MongoDbService[T]) DeleteOne(id any) (int64, error) {
	objectId, _ := primitive.ObjectIDFromHex(id.(string))
	filter := bson.D{{"_id", objectId}}
	result, err := database.Collection(collectionName).DeleteOne(getContext(), filter)

	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, errors.New("no record found")
	}
	return result.DeletedCount, nil
}

func (db *MongoDbService[T]) BulkInsert(documents []T) (T, error) {

	docs := make([]interface{}, len(documents))
	for i, v := range documents {
		docs[i] = v
	}

	result, err := database.Collection(collectionName).InsertMany(getContext(), docs, options.InsertMany())

	if err != nil {
		return documents[0], err
	}
	if len(result.InsertedIDs) != len(documents) {
		return documents[0], errors.New("bulk insertion failed")
	}
	return documents[0], nil
}
