package middleware

import (
	"net/http"
	"strings"
)

// CORS returns middleware that, when allowOrigin is non-empty, sets CORS response headers
// and short-circuits OPTIONS with 204. allowOrigin is typically a single origin or "*".
// When allowOrigin is empty, the next handler is invoked unchanged.
func CORS(allowOrigin string) func(http.Handler) http.Handler {
	ao := strings.TrimSpace(allowOrigin)
	if ao == "" {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Origin", ao)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", ao)
			next.ServeHTTP(w, r)
		})
	}
}
