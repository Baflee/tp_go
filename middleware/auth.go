package middleware

import (
	"net/http"
	"os"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != "Bearer "+os.Getenv("AUTH_TOKEN") {
			httpErrorMsg := "HTTP Error: " + strconv.Itoa(http.StatusUnauthorized) + " - Unauthorized access to the api"
			http.Error(w, httpErrorMsg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
