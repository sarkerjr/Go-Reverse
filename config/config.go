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
	RateLimit       RateLimitConfig
}

type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize        int
	Enabled          bool
}

func LoadConfig() Config {
	return Config{
		Port:            getEnv("PORT"),
		OriginServerURL: getEnv("ORIGIN_SERVER_URL"),
		RateLimit: RateLimitConfig{
			RequestsPerSecond: getEnvAsIntWithDefault("RATE_LIMIT_RPS", 5),
			BurstSize:        getEnvAsIntWithDefault("RATE_LIMIT_BURST", 10),
			Enabled:          getEnvAsBoolWithDefault("RATE_LIMIT_ENABLED", true),
		},
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

func getEnvAsIntWithDefault(name string, defaultValue int) int {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Environment variable %s is not a valid integer, using default value %d", name, defaultValue)
		return defaultValue
	}
	return intValue
}

func getEnvAsBoolWithDefault(name string, defaultValue bool) bool {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Environment variable %s is not a valid boolean, using default value %v", name, defaultValue)
		return defaultValue
	}
	return boolValue
}
