package main

import (
	"log"
	"net/http"
	"os"
	"tp_go/db"
	"tp_go/dictionary"
	"tp_go/middleware"
	"tp_go/router"

	"github.com/joho/godotenv"
)

func main() {
	//Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Establishes a connection to MongoDB.
	// This will block until the connection is established or fails.
	db.ConnectDB(os.Getenv("MONGO_URI"))

	// Initialize important elements of the dictionary
	dictionary.InitializeDictionary()

	// Initializes the router for your HTTP server.
	r := router.InitRouter("db_dictionary_prod")

	// Applies logging middleware to all requests.
	r.Use(middleware.LoggingMiddleware)

	// Applies authentication middleware to all requests.
	r.Use(middleware.AuthMiddleware)

	// Starts the HTTP server on port 8080.
	http.ListenAndServe(":8080", r)
}
