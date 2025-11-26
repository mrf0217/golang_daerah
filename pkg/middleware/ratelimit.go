package middleware

// Request Flow Link:
// main.go wraps every HTTP handler with RateLimitMiddleware, so each incoming request first
// travels through the logic in this file before any handler/service code runs.

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	requests        map[string]*TokenBucket
	mutex           sync.RWMutex
	rate            int           // requests per minute
	burst           int           // maximum burst capacity
	cleanupInterval time.Duration // cleanup interval for old entries
}

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens     int
	lastRefill time.Time
	rate       int
	burst      int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		requests:        make(map[string]*TokenBucket),
		rate:            rate,
		burst:           burst,
		cleanupInterval: 5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.startCleanup()

	return rl
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	bucket, exists := rl.requests[key]

	if !exists {
		// Create new bucket
		bucket = &TokenBucket{
			tokens:     rl.burst - 1, // Allow first request
			lastRefill: now,
			rate:       rl.rate,
			burst:      rl.burst,
		}
		rl.requests[key] = bucket
		return true
	}

	// Refill tokens based on time passed
	timePassed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(timePassed.Minutes()) * bucket.rate

	if tokensToAdd > 0 {
		bucket.tokens = min(bucket.tokens+tokensToAdd, bucket.burst)
		bucket.lastRefill = now
	}

	// Check if we have tokens available
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// startCleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) startCleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, bucket := range rl.requests {
			// Remove entries older than 1 hour
			if now.Sub(bucket.lastRefill) > time.Hour {
				delete(rl.requests, key)
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rate, burst int) func(http.HandlerFunc) http.HandlerFunc {
	limiter := NewRateLimiter(rate, burst)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get client IP as the key for rate limiting
			clientIP := getClientIP(r)

			if !limiter.Allow(clientIP) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{
					"status": false,
					"data": [],
					"message": "Rate limit exceeded. Please try again later."
				}`)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
