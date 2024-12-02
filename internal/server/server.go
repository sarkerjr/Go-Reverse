package server

import (
	"log"
	"net/http"

	"sarkerjr.com/go-reverse/internal/logger"
	"sarkerjr.com/go-reverse/internal/proxy"
	rate_limiter "sarkerjr.com/go-reverse/internal/rate-limiter"
)

func StartServer() error {
	proxySrv := proxy.NewProxy()

	// Create a new IP-based rate limiter instance (5 req/sec, burst of 10)
	ipRateLimiter, err := rate_limiter.NewIPBasedRateLimiter(5, 10)
	if err != nil {
		return err
	}

	var handler http.Handler = http.HandlerFunc(proxySrv.HandleRequest)

	// Apply middlewares in order
	handler = logger.ErrorHandlerMiddleware()(handler)
	handler = ipRateLimiter.Middleware(handler)
	handler = logger.LoggingMiddleware()(handler)

	// Register the handler
	http.Handle("/", handler)

	config := proxy.GetConfig()
	log.Printf("Starting reverse proxy server on port %s...", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
