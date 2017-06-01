package middleware

import "net/http"

// ContentTypeMiddleware wraps the Handler and ensures that
// Content-Type: application/json is set
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
