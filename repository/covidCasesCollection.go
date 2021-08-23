package repository

import (
	"context"
	"crud/config"
	"crud/entity"

	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCovidCasesCollection() *mongo.Collection {
	return config.MongoClient.Database("inshorts").Collection("covid_cases")
}

func UpsertCovidData(sc *mongo.SessionContext, data *entity.CovidCases) (id interface{}, err error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: data.Region}}
	update := bson.D{primitive.E{Key: "$set", Value: data}}
	var res *mongo.UpdateResult
	if sc != nil {
		res, err = getCovidCasesCollection().UpdateOne(*sc, filter, update, opts)
	} else {
		res, err = getCovidCasesCollection().UpdateOne(context.TODO(), filter, update, opts)
	}
	if err != nil {
		log.Println("Exception in upserting doc", err)
		return nil, errors.New("Exception while upserting for state " + data.Region)
	}
	return res.UpsertedID, nil
}

func FindCovidDataByState(region string) (data *entity.CovidCases, err error) {
	var result entity.CovidCases
	filter := bson.D{primitive.E{Key: "_id", Value: region}}
	collection := getCovidCasesCollection()
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
