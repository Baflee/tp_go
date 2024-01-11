package middleware

import (
	"net/http"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != "Bearer RaccoonAreTheBestAnimalsExisting" {
			httpErrorMsg := "HTTP Error: " + strconv.Itoa(http.StatusUnauthorized) + " - Unauthorized access to the api"
			http.Error(w, httpErrorMsg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
