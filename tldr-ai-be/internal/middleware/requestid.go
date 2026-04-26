package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

const xRequestIDHeader = "X-Request-ID"

// RequestID ensures each request has a request id, echoed from the client
// or generated, on the request context and the response X-Request-ID header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get(xRequestIDHeader))
		if id == "" {
			id = generateID()
		}
		r.Header.Set(xRequestIDHeader, id)
		w.Header().Set(xRequestIDHeader, id)
		next.ServeHTTP(w, r)
	})
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
