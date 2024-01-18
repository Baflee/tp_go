package router

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"tp_go/db"
	"tp_go/dictionary"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

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
	dictionary.InitializeDictionary()
}

func TestRouter(t *testing.T) {
	// Initialize the actual router
	router := InitRouter("db_dictionary_test")

	// Log all registered routes (optional)
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, err := route.GetPathTemplate()
		if err == nil {
			t.Logf("Registered Route: %s", template)
		}
		return nil
	})

	// Create a test server using httptest with the router
	server := httptest.NewServer(router)
	defer server.Close()
	t.Logf("Test Server URL: %s", server.URL)

	// Define your test cases
	testCases := []struct {
		name       string
		method     string
		path       string
		statusCode int
		expected   string
	}{
		{
			name:       "TestRouterPathAdd",
			method:     "POST",
			path:       "/dictionary/add?word=testword&definition=testdefinition",
			statusCode: http.StatusOK,
			expected:   "Word testword Added",
		},
		{
			name:       "TestRouterPathGet",
			method:     "GET",
			path:       "/dictionary/testword",
			statusCode: http.StatusOK,
			expected:   "'testword' : 'testdefinition'",
		},
		{
			name:       "TestRouterPathList",
			method:     "GET",
			path:       "/dictionary",
			statusCode: http.StatusOK,
			expected:   "testword: testdefinition",
		},
		{
			name:       "TestRouterPathRemove",
			method:     "DELETE",
			path:       "/dictionary/delete/testword",
			statusCode: http.StatusOK,
			expected:   "Success : Word 'testword' has been removed",
		},
	}
	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create and send the request
			req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check the HTTP status code
			if resp.StatusCode != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, resp.StatusCode)
			}

			// Check the response body (if expected)
			if tc.expected != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %v", err)
				}
				if !strings.Contains(string(body), tc.expected) {
					t.Errorf("Expected response body to contain:\n%s\nGot:\n%s", tc.expected, string(body))
				}
			}
		})
	}
}
