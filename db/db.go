package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client         *mongo.Client
	TestClient     *mongo.Client
	TestCollection *mongo.Collection
)

func ConnectDB(uri string) {
	// Attempt to connect
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Db Error Connecting : '%s'", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Db Error Checking Connection : '%s'", err)
	}

	Client = client
	log.Println("Successful connection to MongoDB.")
}
