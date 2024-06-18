package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port            string
	OriginServerURL string
	AuthToken       string
	RateLimit       int
}

func LoadConfig() Config {
	return Config{
		Port:            getEnv("PORT"),
		OriginServerURL: getEnv("ORIGIN_SERVER_URL"),
	}
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func getEnvAsInt(name string) int {
	value, err := strconv.Atoi(getEnv(name))
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer", name)
	}
	return value
}
