package DB

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
	collectionName = "task"
)

func DBInstance() *mongo.Client {

	MongoServerURL := os.ExpandEnv("mongodb://${DB_USERNAME}:${DB_PASSWORD}@$DB_HOST:$DB_PORT/$DB_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoServerURL))
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
	var collection *mongo.Collection = client.Database("taskmanager").Collection(collectionName)
	return collection
}

func EstablishConnection() {
	Collection = OpenCollection(DBInstance(), collectionName)
	if Collection == nil {
		log.Fatal("No Collection found")
	}
}
