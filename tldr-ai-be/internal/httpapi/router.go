package httpapi

import (
	"net/http"

	"tldr-ai-be/internal/middleware"
	"tldr-ai-be/internal/ratelimit"
)

// NewHandler registers /health and /api/processText and wraps the mux with
// Recover(SecurityHeaders(RequestID(CORS(mux)))).
// trustProxy and corsAllow are passed to client IP and CORS. limiter can be nil (no limit).
func NewHandler(d *RouterDeps, trustProxy bool, corsAllow string, limiter *ratelimit.Limiter) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	pt := d.processText
	if limiter != nil {
		p := limiter.HTTPMiddleware(func(r *http.Request) string { return ClientIP(r, trustProxy) })
		mux.Handle("POST /api/processText", p(http.HandlerFunc(pt)))
	} else {
		mux.Handle("POST /api/processText", http.HandlerFunc(pt))
	}
	h := http.Handler(mux)
	h = middleware.CORS(corsAllow)(h)
	h = middleware.RequestID(h)
	h = middleware.SecurityHeaders(h)
	h = middleware.Recover(h)
	return h
}
