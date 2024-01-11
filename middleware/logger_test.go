package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogMiddleware(t *testing.T) {
	// Create a temporary log file for testing
	logFileName := "requests.log"

	// Initialize the logger with the temporary log file
	log.SetOutput(os.Stdout) // Redirect logger output to the console for testing

	// Create a handler that will be wrapped by the LoggingMiddleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with the LoggingMiddleware
	ts := httptest.NewServer(LoggingMiddleware(handler))
	defer ts.Close()

	// Send a test request to the server
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/dictionary", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Ensure that the log file contains the expected log message
	logContents, err := os.ReadFile(logFileName)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	expectedLogMessage := "Method: GET, Path: /dictionary, Status: 200"
	if !strings.Contains(string(logContents), expectedLogMessage) {
		t.Errorf("Log file does not contain expected log message: %s", expectedLogMessage)
	}

	// Remove the log file after reading it
	if err := os.Remove(logFileName); err != nil {
		t.Fatalf("Failed to remove log file: %v", err)
	}
}
