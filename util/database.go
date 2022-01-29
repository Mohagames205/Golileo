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

func InitFs() error {
	workingDirectory, _ := os.Getwd()
	err := os.MkdirAll(workingDirectory+"/images/", os.ModeDir)

	return err
}

func InitDatabase() {
	/*
	   Connect to the cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

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
