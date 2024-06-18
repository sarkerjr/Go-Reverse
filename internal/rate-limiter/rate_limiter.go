package rate_limiter

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	rate      int // number of tokens generated per second
	burst     int // maximum number of tokens
	tokens    int // number of tokens available
	mutex     sync.Mutex
	timestamp time.Time
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
	return &RateLimiter{
		rate:      rate,
		burst:     burst,
		tokens:    burst,
		timestamp: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.timestamp).Seconds()
	rl.timestamp = now

	// refill the bucket
	rl.tokens += int(float64(rl.rate) * elapsed)
	if rl.tokens > rl.burst {
		rl.tokens = rl.burst
	}

	// Check if there are enough tokens
	if rl.tokens > 0 {
		rl.tokens--
		log.Printf("Allowing request. Tokens left: %d", rl.tokens)
		return true
	}

	log.Println("Too many requests. Rate limit exceeded.")
	return false
}

type IPBasedRateLimiter struct {
	limiters map[string]*RateLimiter
	rate     int
	burst    int
	mutex    sync.Mutex
}

func NewIPBasedRateLimiter(rate int, burst int) *IPBasedRateLimiter {
	return &IPBasedRateLimiter{
		limiters: make(map[string]*RateLimiter),
		rate:     rate,
		burst:    burst,
	}
}

func (iprl *IPBasedRateLimiter) getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// getRateLimiter returns the rate limiter for the given IP address
func (iprl *IPBasedRateLimiter) getRateLimiter(ip string) *RateLimiter {
	iprl.mutex.Lock()
	defer iprl.mutex.Unlock()

	if limiter, exists := iprl.limiters[ip]; exists {
		return limiter
	}

	limiter := NewRateLimiter(iprl.rate, iprl.burst)
	iprl.limiters[ip] = limiter
	return limiter
}

// Middleware function
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
