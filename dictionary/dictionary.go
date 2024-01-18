package dictionary

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tp_go/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeDictionary() {
	go handleAddRequests()
	go handleRemoveRequests()
	go handleGetRequests()
	go handleListRequests()
}

var (
	addChan      = make(chan AddRequest)
	removeChan   = make(chan GetRemoveRequest)
	getChan      = make(chan GetRemoveRequest)
	listChan     = make(chan ListRequest)
	responseChan = make(chan Response)
)

var collection *mongo.Collection

type AddRequest struct {
	Db         string
	Collection string
	Key        string
	Value      string
	Response   chan Response
}

type GetRemoveRequest struct {
	Db         string
	Collection string
	Key        string
	Response   chan Response
}

type ListRequest struct {
	Db         string
	Collection string
	Response   chan Response
}

type Response struct {
	Result string
	Err    error
	Http   int
}

func Add(db string, collection string, key string, value string) (string, error, int) {
	responseChan := make(chan Response)
	addChan <- AddRequest{Db: db, Collection: collection, Key: key, Value: value, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func Remove(db string, collection string, key string) (string, error, int) {
	responseChan := make(chan Response)
	removeChan <- GetRemoveRequest{Db: db, Collection: collection, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func Get(db string, collection string, key string) (string, error, int) {
	responseChan := make(chan Response)
	getChan <- GetRemoveRequest{Db: db, Collection: collection, Key: key, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func List(db string, collection string) (string, error, int) {
	responseChan := make(chan Response)
	listChan <- ListRequest{Db: db, Collection: collection, Response: responseChan}
	response := <-responseChan
	return response.Result, response.Err, response.Http
}

func handleAddRequests() {
	for req := range addChan {
		collection = db.Client.Database(req.Db).Collection(req.Collection)

		// Setup the filter to check if the word already exists
		filter := bson.M{"key": req.Key}
		var result bson.M
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Help to avoid being stuck indefinitely in this request
		defer cancel()
		err := collection.FindOne(ctx, filter).Decode(&result)

		// If the word already exists, return an error
		if err == nil {
			req.Response <- Response{"", fmt.Errorf("Word '%s' already exists", req.Key), http.StatusConflict}
			continue
		}

		// Otherwise, add the new word
		doc := bson.M{"key": req.Key, "value": req.Value}
		_, err = collection.InsertOne(ctx, doc)
		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}

		// Send confirmation of word addition
		req.Response <- Response{fmt.Sprintf("Success : Word '%s' has been added", req.Key), nil, http.StatusOK}
	}
}

func handleRemoveRequests() {
	for req := range removeChan {
		collection = db.Client.Database(req.Db).Collection(req.Collection)

		// Attempt to remove the word
		filter := bson.M{"key": req.Key}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		result, err := collection.DeleteOne(ctx, filter)

		// Manage errors or missing words to be deleted
		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}
		if result.DeletedCount == 0 {
			req.Response <- Response{"", fmt.Errorf("Word '%s' does not exist", req.Key), http.StatusNotFound}
			continue
		}

		// Send confirmation of word removal
		req.Response <- Response{fmt.Sprintf("Success : Word '%s' has been removed", req.Key), nil, http.StatusOK}
	}
}

func handleGetRequests() {
	for req := range getChan {
		collection = db.Client.Database(req.Db).Collection(req.Collection)

		// Setup the filter to find the document with the given key
		filter := bson.M{"key": req.Key}
		var result bson.M
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := collection.FindOne(ctx, filter).Decode(&result)

		// Handle errors or missing words
		if err != nil {
			req.Response <- Response{"", fmt.Errorf("Word '%s' not found", req.Key), http.StatusNotFound}
			continue
		}

		value, ok := result["value"].(string)
		if !ok {
			req.Response <- Response{"", fmt.Errorf("Invalid data format for key: %s", req.Key), http.StatusInternalServerError}
			continue
		}

		// Send the formatted string of the searched word and its value
		req.Response <- Response{fmt.Sprintf("'%s' : '%s'", req.Key, value), nil, http.StatusOK}
	}
}

func handleListRequests() {
	for req := range listChan {
		collection = db.Client.Database(req.Db).Collection(req.Collection)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Retrieve all documents from the collection
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			req.Response <- Response{"", err, http.StatusInternalServerError}
			continue
		}

		// Format the results as a string for the response
		var list strings.Builder
		for _, result := range results {
			key := result["key"].(string)
			value := result["value"].(string)
			list.WriteString(key + ": " + value + "\n")
		}

		// Send the formatted string of all words and values
		req.Response <- Response{list.String(), nil, http.StatusOK}
	}
}
