package middleware

import (
	"net/http"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform authentication checks

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
