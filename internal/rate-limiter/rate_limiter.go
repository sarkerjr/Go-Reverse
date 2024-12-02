package rate_limiter

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultCleanupInterval = 5 * time.Minute
)

// RateLimiter implements a token bucket algorithm for rate limiting.
type RateLimiter struct {
	rate      int       // number of tokens generated per second
	burst     int       // maximum number of tokens
	tokens    int       // current number of tokens available
	timestamp time.Time // last time tokens were generated
	mx        sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified rate and burst size.
func NewRateLimiter(rate, burst int) (*RateLimiter, error) {
	if rate <= 0 || burst <= 0 {
		return nil, fmt.Errorf("rate and burst must be positive values")
	}

	return &RateLimiter{
		rate:      rate,
		burst:     burst,
		tokens:    burst,
		timestamp: time.Now(),
	}, nil
}

// Allow checks if a request should be allowed based on the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.mx.Lock()
	defer rl.mx.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.timestamp).Seconds()
	rl.timestamp = now

	// Refill the bucket
	rl.tokens += int(float64(rl.rate) * elapsed)
	if rl.tokens > rl.burst {
		rl.tokens = rl.burst
	}

	if rl.tokens > 0 {
		rl.tokens--
		log.Printf("Rate limit status: %d tokens remaining", rl.tokens)
		return true
	}

	log.Printf("Rate limit exceeded: 0 tokens remaining")
	return false
}

// IPBasedRateLimiter manages rate limiters for different IP addresses.
type IPBasedRateLimiter struct {
	limiters map[string]*RateLimiter
	rate     int
	burst    int
	mutex    sync.RWMutex
}

// NewIPBasedRateLimiter creates a new IP-based rate limiter.
func NewIPBasedRateLimiter(rate, burst int) (*IPBasedRateLimiter, error) {
	if rate <= 0 || burst <= 0 {
		return nil, fmt.Errorf("rate and burst must be positive values")
	}

	limiter := &IPBasedRateLimiter{
		limiters: make(map[string]*RateLimiter),
		rate:     rate,
		burst:    burst,
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter, nil
}

// getClientIP extracts the client IP address from the request.
func (iprl *IPBasedRateLimiter) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// getRateLimiter returns the rate limiter for the given IP address.
func (iprl *IPBasedRateLimiter) getRateLimiter(ip string) *RateLimiter {
	iprl.mutex.RLock()
	limiter, exists := iprl.limiters[ip]
	iprl.mutex.RUnlock()

	if exists {
		return limiter
	}

	iprl.mutex.Lock()
	defer iprl.mutex.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists = iprl.limiters[ip]; exists {
		return limiter
	}

	// Create a new rate limiter
	limiter, _ = NewRateLimiter(iprl.rate, iprl.burst)
	iprl.limiters[ip] = limiter
	return limiter
}

// cleanup periodically removes expired rate limiters to prevent memory leaks.
func (iprl *IPBasedRateLimiter) cleanup() {
	ticker := time.NewTicker(defaultCleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		iprl.mutex.Lock()
		for ip, limiter := range iprl.limiters {
			if time.Since(limiter.timestamp) > defaultCleanupInterval {
				delete(iprl.limiters, ip)
			}
		}
		iprl.mutex.Unlock()
	}
}

// Middleware returns an http.Handler middleware that applies rate limiting.
func (iprl *IPBasedRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := iprl.getClientIP(r)
		limiter := iprl.getRateLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
