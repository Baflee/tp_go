package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ouvrir le fichier journal en mode append
		file, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Enregistrer la requête
		log.SetOutput(file)
		log.Printf("Time: %v, Method: %s, Path: %s\n", time.Now(), r.Method, r.URL.Path)

		// Passer au prochain handler
		next.ServeHTTP(w, r)
	})
}
