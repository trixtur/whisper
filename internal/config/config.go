package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	OpenAIAPIKey string
	Model        string
	Language     string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	cfg := &Config{
		OpenAIAPIKey: apiKey,
		Model:        getEnvOrDefault("STT_MODEL", "whisper-1"),
		Language:     getEnvOrDefault("STT_LANGUAGE", "en"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
