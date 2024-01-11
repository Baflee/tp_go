package dictionary

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictionaryAdd(t *testing.T) {
	filePathTest := "dictionary_test.txt"
	file, _ := os.Create(filePathTest)

	// Success Case
	result, err, statusCode := Add(filePathTest, "new_word", "definition")
	assert.NoError(t, err)
	assert.Equal(t, "Success", result)
	assert.Equal(t, 200, statusCode)

	// Key Already Exists Case
	Add(filePathTest, "existing_key", "definition")
	result, err, statusCode = Add("dictionary_test.txt", "existing_key", "definition")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Equal(t, 409, statusCode)

	// Invalid File Path Case
	result, err, statusCode = Add("invalid_file_path.txt", "new_word", "definition")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File 'invalid_file_path.txt' not found")
	assert.Equal(t, 404, statusCode)

	defer func() {
		// Close the file before removing it
		file.Close()
		// Remove the temporary file
		if err := os.Remove(filePathTest); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()
}

func TestDictionaryRemove(t *testing.T) {
	filePathTest := "dictionary_test.txt"
	file, _ := os.Create(filePathTest)
	Add(filePathTest, "existing_key", "definition")

	// Success Case
	result, err, statusCode := Remove(filePathTest, "existing_key")
	assert.NoError(t, err)
	assert.Equal(t, "Success", result)
	assert.Equal(t, 200, statusCode)

	// Key Not Found Case
	result, err, statusCode = Remove(filePathTest, "non_existing_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	assert.Equal(t, 404, statusCode)

	// Invalid File Path Case
	result, err, statusCode = Remove("invalid_file_path.txt", "existing_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File 'invalid_file_path.txt' not found")
	assert.Equal(t, 404, statusCode)

	defer func() {
		// Close the file before removing it
		file.Close()
		// Remove the temporary file
		if err := os.Remove(filePathTest); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()
}

func TestDictionaryGet(t *testing.T) {
	filePathTest := "dictionary_test.txt"
	file, _ := os.Create(filePathTest)
	Add(filePathTest, "existing_key", "definition")

	// Success Case
	result, err, statusCode := Get(filePathTest, "existing_key")
	assert.NoError(t, err)
	assert.Contains(t, result, "")
	assert.Equal(t, 200, statusCode)

	// Key Not Found Case
	result, err, statusCode = Get(filePathTest, "non_existing_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, 404, statusCode)

	// Invalid File Path Case
	result, err, statusCode = Get("invalid_file_path.txt", "existing_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File 'invalid_file_path.txt' not found")
	assert.Equal(t, 404, statusCode)

	defer func() {
		// Close the file before removing it
		file.Close()
		// Remove the temporary file
		if err := os.Remove(filePathTest); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()
}

func TestDictionaryList(t *testing.T) {
	filePathTest := "dictionary_test.txt"
	file, _ := os.Create(filePathTest)

	Add(filePathTest, "word1", "definition1")
	Add(filePathTest, "word2", "definition2")

	// Success Case with Multiple Words
	result, err, statusCode := List(filePathTest)
	assert.NoError(t, err)
	assert.Contains(t, result, "word1:definition1")
	assert.Contains(t, result, "word2:definition2")
	assert.Equal(t, 200, statusCode)

	Remove(filePathTest, "word1")
	Remove(filePathTest, "word2")

	// Empty File Case
	result, err, statusCode = List(filePathTest)
	assert.Equal(t, "Empty", result)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	// Invalid File Path Case
	result, err, statusCode = List("invalid_file_path.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File 'invalid_file_path.txt' not found")
	assert.Equal(t, 404, statusCode)

	defer func() {
		// Close the file before removing it
		file.Close()
		// Remove the temporary file
		if err := os.Remove(filePathTest); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()
}
