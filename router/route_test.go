package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func TestRouter(t *testing.T) {
	filePathTest := "dictionary_test.txt"
	file, _ := os.Create(filePathTest)

	router := InitRouter(filePathTest)
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
	t.Logf(server.URL)

	testCases := []struct {
		name       string
		method     string
		path       string
		statusCode int
		expected   string
	}{
		{
			name:       "Add a Word and definition",
			method:     "POST",
			path:       "/dictionary/add?word=testword&definition=testdefinition",
			statusCode: http.StatusOK,
			expected:   "Word testword Added",
		},
		{
			name:       "Define a Word",
			method:     "GET",
			path:       "/dictionary/testword",
			statusCode: http.StatusOK,
			expected:   "Definition of testword: testdefinition",
		},
		{
			name:       "List every Words with their definitions",
			method:     "GET",
			path:       "/dictionary",
			statusCode: http.StatusOK,
			expected:   "Word Lists: \ntestword:testdefinition",
		},
		{
			name:       "Remove a Word",
			method:     "DELETE",
			path:       "/dictionary/delete/testword",
			statusCode: http.StatusOK,
			expected:   "Word testword Deleted",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Send the request
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
				if string(body) != tc.expected {
					t.Errorf("Expected response body:\n%s\nGot:\n%s", tc.expected, string(body))
				}
			}
		})
	}

	defer func() {
		// Close the file before removing it
		file.Close()
		// Remove the temporary file
		if err := os.Remove(filePathTest); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()
}
