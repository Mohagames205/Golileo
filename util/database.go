package util

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var Client *mongo.Client

func Database() *mongo.Database {
	return Client.Database("galileo")
}

func InitDatabase() {
	/*
	   Connect to the cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT")))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	Client = client

}
