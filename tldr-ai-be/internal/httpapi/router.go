package httpapi

import (
	"net/http"
	"strings"

	"tldr-ai-be/internal/middleware"
	"tldr-ai-be/internal/ratelimit"
)

// NewHandler registers routes and wraps the mux with
// Recover(SecurityHeaders(RequestID(CORS(mux)))).
func NewHandler(d *RouterDeps, trustProxy bool, corsAllow string, limiter *ratelimit.Limiter) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.Handle("GET /api/usage", http.HandlerFunc(d.usageGet))
	if strings.TrimSpace(d.UsageResetSecret) != "" {
		mux.Handle("POST /api/admin/usage-reset", http.HandlerFunc(d.usageAdminReset))
	}
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
