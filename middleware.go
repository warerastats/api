package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/warerastats/api/graph"
)

// writeGraphQLError responds with a GraphQL-shaped JSON error body so GraphQL
// clients can parse it instead of choking on plain text.
func writeGraphQLError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"errors": []map[string]any{{"message": msg}},
	})
}

// splitEnv reads a comma-separated environment variable into a trimmed,
// non-empty slice of values.
func splitEnv(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// --- API key authentication ---

// apiKeyAuth rejects any request that does not present a known API key via the
// `Authorization: Bearer <key>` header. Keys are read once from the API_KEYS
// environment variable (comma separated). The same key is used by the SvelteKit
// frontend and any future third-party consumer.
func apiKeyAuth(next http.Handler) http.Handler {
	keys := make(map[string]struct{})
	for _, k := range splitEnv("API_KEYS") {
		keys[k] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		if _, ok := keys[bearerToken(r)]; !ok {
			writeGraphQLError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// bearerToken extracts the token from an `Authorization: Bearer <token>` header.
func bearerToken(r *http.Request) string {
	const prefix = "Bearer "
	h := r.Header.Get("Authorization")
	if len(h) > len(prefix) && strings.EqualFold(h[:len(prefix)], prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	return ""
}

// --- CORS ---

// cors echoes back the request Origin only when it appears in the CORS_ORIGINS
// allow-list and answers preflight requests. Browsers enforce the policy; this
// is not an authentication mechanism.
func cors(next http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, o := range splitEnv("CORS_ORIGINS") {
		allowed[o] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				h := w.Header()
				h.Set("Access-Control-Allow-Origin", origin)
				h.Set("Vary", "Origin")
				h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				h.Set("Access-Control-Max-Age", "86400")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- per-IP rate limiting ---

// ipRateLimiter throttles the public playground query endpoint to one request
// per interval per client IP.
type ipRateLimiter struct {
	mu      sync.Mutex
	clients map[string]*rateClient
	limit   rate.Limit
	burst   int
	idleFor time.Duration
}

type rateClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// newIPRateLimiter builds a limiter allowing one request every `every` duration
// per IP, evicting clients that have been idle for longer than 10 minutes.
func newIPRateLimiter(every time.Duration) *ipRateLimiter {
	l := &ipRateLimiter{
		clients: make(map[string]*rateClient),
		limit:   rate.Every(every),
		burst:   1,
		idleFor: 10 * time.Minute,
	}
	go l.cleanupLoop()
	return l
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	c, ok := l.clients[ip]
	if !ok {
		c = &rateClient{limiter: rate.NewLimiter(l.limit, l.burst)}
		l.clients[ip] = c
	}
	c.lastSeen = time.Now()
	return c.limiter.Allow()
}

func (l *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		for ip, c := range l.clients {
			if time.Since(c.lastSeen) > l.idleFor {
				delete(l.clients, ip)
			}
		}
		l.mu.Unlock()
	}
}

// middleware wraps next with the per-IP rate limit, returning 429 when exceeded.
func (l *ipRateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		if !l.allow(clientIP(r)) {
			writeGraphQLError(w, http.StatusTooManyRequests, "rate limit exceeded (1 req per 10 seconds): try again later")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP resolves the caller IP, preferring the first X-Forwarded-For entry
// when running behind a reverse proxy.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// --- playground context ---

// playgroundContext tags every request to the public query endpoint so that
// resolvers apply the public time-window caps.
func playgroundContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(graph.WithPlayground(r.Context())))
	})
}
