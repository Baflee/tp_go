package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	//Load the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestConnectDB(t *testing.T) {
	// Call ConnectDB
	ConnectDB(os.Getenv("MONGO_URI"))

	// Check if the client is not nil
	if Client == nil {
		t.Fatal("Expected a non-nil MongoDB client, got nil")
	}

	// Optionally, disconnect after the test
	err := Client.Disconnect(context.TODO())
	if err != nil {
		t.Fatalf("Failed to disconnect from MongoDB: %s", err)
	}
}
