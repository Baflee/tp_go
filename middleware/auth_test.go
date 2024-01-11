package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a handler that will be wrapped by the AuthMiddleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with the AuthMiddleware
	ts := httptest.NewServer(AuthMiddleware(handler))
	defer ts.Close()

	// Test cases with different authorization tokens
	testCases := []struct {
		name          string
		authorization string
		expectedCode  int
	}{
		{
			name:          "Valid Token",
			authorization: "Bearer RaccoonAreTheBestAnimalsExisting",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Invalid Token",
			authorization: "Bearer InvalidToken",
			expectedCode:  http.StatusUnauthorized,
		},
		{
			name:          "No Token",
			authorization: "",
			expectedCode:  http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Authorization", tc.authorization)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, resp.StatusCode)
			}
		})
	}
}
