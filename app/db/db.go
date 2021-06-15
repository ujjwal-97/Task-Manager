package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Collection *mongo.Collection
)

const (
	CollectionName = "task"
)

func DBInstance() *mongo.Client {

	MongoServerURL := os.ExpandEnv("mongodb://$DB_HOST:$DB_PORT/$DB_NAME")
	credential := options.Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoServerURL).SetAuth(credential))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to MongoDb")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(os.Getenv("DB_NAME")).Collection(collectionName)
	return collection
}

func EstablishConnection() {
	Collection = OpenCollection(DBInstance(), CollectionName)
	if Collection == nil {
		log.Fatal("No Collection found")
	}
}
