package httpapi

import (
	"net"
	"net/http"
	"strings"
)

// ClientIP returns a client key suitable for per-IP limits and logs.
// When trustProxy is true, the first non-empty value in X-Forwarded-For is used (comma-separated, left to right).
// Otherwise, the host part of RemoteAddr is used.
func ClientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
		if xff != "" {
			for _, p := range strings.Split(xff, ",") {
				ip := strings.TrimSpace(p)
				if ip != "" {
					return ip
				}
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
