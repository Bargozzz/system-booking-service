package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// bucket implements a simple token-bucket per client.
type bucket struct {
	tokens     float64
	lastRefill time.Time
}

// Limiter is a per-IP token-bucket rate limiter safe for concurrent use.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64 // tokens added per second
	capacity float64 // max tokens (burst size)
}

// New creates a limiter allowing `ratePerSecond` sustained requests per
// second per client IP, with a burst capacity of `burst`.
func New(ratePerSecond float64, burst float64) *Limiter {
	return &Limiter{
		buckets:  make(map[string]*bucket),
		rate:     ratePerSecond,
		capacity: burst,
	}
}

func (l *Limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[key]
	now := time.Now()
	if !ok {
		b = &bucket{tokens: l.capacity - 1, lastRefill: now}
		l.buckets[key] = b
		return true
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.lastRefill = now

	if b.tokens < 1 {
		return false
	}
	b.tokens -= 1
	return true
}

// Middleware wraps a handler, rejecting requests over the limit with 429.
func (l *Limiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		if !l.allow(host) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded, please slow down"}`))
			return
		}
		next.ServeHTTP(w, r)
	}
}
