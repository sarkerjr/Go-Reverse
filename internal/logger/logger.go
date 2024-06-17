package logger

import (
	"log"
	"net/http"
)

type HandlerDecorator func(http.Handler) http.Handler

// LoggingMiddleware logs the incoming request and the outgoing response
func LoggingMiddleware() HandlerDecorator {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			log.Printf("Completed request: %s %s", r.Method, r.URL.Path)
		})
	}
}

// ErrorHandlerMiddleware recovers from panics and returns a 500 Internal Server Error
func ErrorHandlerMiddleware() HandlerDecorator {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered from panic: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
