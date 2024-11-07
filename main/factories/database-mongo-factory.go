package factories

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db_mongo_con *mongo.Client

func NewDatabaseMongoOpenConnection() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return errors.New("MONGO_URI is not defined")
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return errors.New(err.Error())
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return errors.New(err.Error())
	}

	db_mongo_con = client
	return nil
}

func NewCloseDatabaseMongoConnection() error {
	return db_mongo_con.Disconnect(context.Background())
}