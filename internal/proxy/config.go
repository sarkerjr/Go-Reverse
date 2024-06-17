// internal/proxy/config.go
package proxy

import (
	"log"
	"os"
	"sync"
)

type Config struct {
	Port      string
	OriginURL string
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("PORT environment variable is not set")
		}

		originURL := os.Getenv("ORIGIN_SERVER_URL")
		if originURL == "" {
			log.Fatal("ORIGIN_SERVER_URL environment variable is not set")
		}

		instance = &Config{
			Port:      port,
			OriginURL: originURL,
		}
	})
	return instance
}
