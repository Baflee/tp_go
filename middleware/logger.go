package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

var file *os.File

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{w, http.StatusOK}
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Open the log file for each request
		file, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening log file: %v", err)
		}
		defer file.Close() // Ensure the file is closed when the request is done

		recorder := newResponseRecorder(w)
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		// Use the opened log file for logging
		log.SetOutput(file)
		log.Printf("Time: %v, Method: %s, Path: %s, Status: %d, Duration: %v\n", start, r.Method, r.URL.Path, recorder.statusCode, duration)
	})
}
