package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != "RaccoonAreTheBestAnimalsExisting" {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
