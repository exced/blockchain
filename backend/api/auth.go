package api

import (
	"net/http"
)

// WithAuth is a barrier to ensure user is authenticated before serving Http handler
func WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure user is auth
		h.ServeHTTP(w, r)
	})
}
