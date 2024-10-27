package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func logFailure(key string, sources []string) {
	logger := NewLogger()
	logger.Errorw("Failed to get key-value pair", "failed sources", sources, "key", key)
}

func GetEnv(key string) string {
	var failedSources []string

	// Attempt to load the .env file only once
	if err := godotenv.Load(); err != nil {
		logger := NewLogger()
		logger.Warn("Could not load .env file", "error", err)
	}

	// Try to get from .env first, then from OS
	if value := getDotEnv(key, &failedSources); value != "" {
		return value
	}

	logFailure(key, failedSources)
	return ""
}

func getDotEnv(key string, failedSources *[]string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	*failedSources = append(*failedSources, ".env")
	return ""
}
