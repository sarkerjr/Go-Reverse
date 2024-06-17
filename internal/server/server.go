package server

import (
	"log"
	"net/http"

	"sarkerjr.com/go-reverse/internal/logger"
	"sarkerjr.com/go-reverse/internal/proxy"
)

func StartServer() {
	proxySrv := proxy.NewProxy()
	http.Handle("/", logger.LoggingMiddleware()(logger.ErrorHandlerMiddleware()(http.HandlerFunc(proxySrv.HandleRequest))))

	config := proxy.GetConfig()
	log.Printf("Starting reverse proxy server on port %s...", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
