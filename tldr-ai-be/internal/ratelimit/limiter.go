package ratelimit

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"tldr-ai-be/internal/errs"
	"tldr-ai-be/internal/web"
)

// Limiter is a per-IP token bucket. RPS==0 means disabled (all requests allowed).
type Limiter struct {
	mu     sync.Mutex
	rps    float64
	burst  float64
	maxIPs int
	ips    map[string]*ipBucket
	fifo   []string
}

type ipBucket struct {
	tokens   float64
	lastUnix int64
}

// Init loads RATE_LIMIT_RPS (0 disables) and RATE_LIMIT_MAX_IPS (default 10_000, floor 1).
// get should return raw env, e.g. os.Getenv. Empty RATE_LIMIT_RPS is treated as 0.
func Init(get func(key string) string) *Limiter {
	s := strings.TrimSpace(get("RATE_LIMIT_RPS"))
	if s == "" {
		return newLimiter(0, defaultMaxIPs(get))
	}
	rps, err := strconv.ParseFloat(s, 64)
	if err != nil || rps < 0 {
		rps = 0
	}
	return newLimiter(rps, defaultMaxIPs(get))
}

func defaultMaxIPs(get func(string) string) int {
	s := strings.TrimSpace(get("RATE_LIMIT_MAX_IPS"))
	if s == "" {
		return 10_000
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 10_000
	}
	return n
}

func newLimiter(rps float64, maxIPs int) *Limiter {
	if maxIPs < 1 {
		maxIPs = 1
	}
	b := rps
	if b < 1 {
		b = 1
	}
	if rps == 0 {
		return &Limiter{rps: 0, maxIPs: maxIPs, ips: make(map[string]*ipBucket), fifo: make([]string, 0)}
	}
	return &Limiter{
		rps:    rps,
		burst:  b * 2,
		maxIPs: maxIPs,
		ips:    make(map[string]*ipBucket),
		fifo:   make([]string, 0, 64),
	}
}

// Allow reports whether a request for ip should be accepted (HTTP 2xx path).
func (l *Limiter) Allow(ip string) bool {
	if l == nil || l.rps == 0 {
		return true
	}
	now := time.Now().UnixNano()
	l.mu.Lock()
	defer l.mu.Unlock()

	if b, ok := l.ips[ip]; ok {
		return l.refillAndTake(b, now)
	}
	if len(l.ips) >= l.maxIPs {
		old := l.fifo[0]
		l.fifo = l.fifo[1:]
		delete(l.ips, old)
	}
	b := &ipBucket{tokens: l.burst, lastUnix: now}
	l.ips[ip] = b
	l.fifo = append(l.fifo, ip)
	return l.refillAndTake(b, now)
}

func (l *Limiter) refillAndTake(b *ipBucket, now int64) bool {
	elapsed := float64(now-b.lastUnix) / 1e9
	if elapsed > 0 {
		b.tokens += elapsed * l.rps
		if b.tokens > l.burst {
			b.tokens = l.burst
		}
		b.lastUnix = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens -= 1
	return true
}

// HTTPMiddleware enforces the limiter; clientIP is called on each request (e.g. return httpapi.ClientIP(r, trustProxy)).
// When disabled (RPS=0) or l is nil, the middleware is a no-op.
func (l *Limiter) HTTPMiddleware(clientIP func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if l == nil || l.rps == 0 {
				next.ServeHTTP(w, r)
				return
			}
			ip := clientIP(r)
			if !l.Allow(ip) {
				web.HandleError(w, r, errs.TooManyRequests(""))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
