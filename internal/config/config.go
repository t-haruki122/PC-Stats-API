package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port        string
	IntervalSec int
	HistorySize int
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		IntervalSec: getEnvInt("INTERVAL_SEC", 30),
		HistorySize: getEnvInt("HISTORY_SIZE", 720),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
