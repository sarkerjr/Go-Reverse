package server

import (
	"log"
	"net/http"

	"sarkerjr.com/go-reverse/internal/logger"
	"sarkerjr.com/go-reverse/internal/proxy"
	rate_limiter "sarkerjr.com/go-reverse/internal/rate-limiter"
)

func StartServer() {
	proxySrv := proxy.NewProxy()

	// create a new IP-based rate limiter instance
	ipRateLimiter := rate_limiter.NewIPBasedRateLimiter(5, 10)

	http.Handle("/",
		logger.LoggingMiddleware()(
			logger.ErrorHandlerMiddleware()(
				ipRateLimiter.Middleware(
					http.HandlerFunc(proxySrv.HandleRequest),
				),
			),
		),
	)

	config := proxy.GetConfig()
	log.Printf("Starting reverse proxy server on port %s...", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
