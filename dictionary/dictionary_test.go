package dictionary

import (
	"log"
	"net/http"
	"os"
	"testing"
	"tp_go/db"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var testClient *mongo.Client
var testCollection *mongo.Collection

var baseTestDB = "db_dictionary_test"
var baseTestCollection = "dictionary"

func init() {
	//Load the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Establishes a connection to MongoDB.
	// This will block until the connection is established or fails.
	db.ConnectDB(os.Getenv("MONGO_URI"))

	// Initialize important elements of the dictionary
	InitializeDictionary()
}

func TestDictionaryAdd(t *testing.T) {

	// Success Case
	result, err, statusCode := Add(baseTestDB, baseTestCollection, "new_word", "definition")
	assert.NoError(t, err)
	assert.Equal(t, "Success : Word 'new_word' has been added", result)
	assert.Equal(t, http.StatusOK, statusCode)

	// Key Already Exists Case
	Add(baseTestDB, baseTestCollection, "existing_key", "definition")
	result, err, statusCode = Add(baseTestDB, baseTestCollection, "existing_key", "definition")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Equal(t, http.StatusConflict, statusCode)

	Remove(baseTestDB, baseTestCollection, "new_word")
	Remove(baseTestDB, baseTestCollection, "existing_key")
}

func TestDictionaryRemove(t *testing.T) {

	Add(baseTestDB, baseTestCollection, "existing_key", "definition")

	// Success Case
	result, err, statusCode := Remove(baseTestDB, baseTestCollection, "existing_key")
	assert.NoError(t, err)
	assert.Equal(t, "Success : Word 'existing_key' has been removed", result)
	assert.Equal(t, http.StatusOK, statusCode)

	// Key Not Found Case
	result, err, statusCode = Remove(baseTestDB, baseTestCollection, "non_existing_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestDictionaryGet(t *testing.T) {

	Add(baseTestDB, baseTestCollection, "existing_key", "definition")

	// Success Case
	result, err, statusCode := Get(baseTestDB, baseTestCollection, "existing_key")
	assert.NoError(t, err)
	assert.Equal(t, "'existing_key' : 'definition'", result)
	assert.Equal(t, http.StatusOK, statusCode)

	// Key Not Found Case
	result, err, statusCode = Get(baseTestDB, baseTestCollection, "non_existing_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, http.StatusNotFound, statusCode)

	Remove(baseTestDB, baseTestCollection, "existing_key")
}

func TestDictionaryList(t *testing.T) {

	Add(baseTestDB, baseTestCollection, "word1", "definition1")
	Add(baseTestDB, baseTestCollection, "word2", "definition2")

	// Success Case with Multiple Words
	result, err, statusCode := List(baseTestDB, baseTestCollection)
	assert.NoError(t, err)
	assert.Contains(t, result, "word1: definition1")
	assert.Contains(t, result, "word2: definition2")
	assert.Equal(t, http.StatusOK, statusCode)

	Remove(baseTestDB, baseTestCollection, "word1")
	Remove(baseTestDB, baseTestCollection, "word2")
}
