package config

import (
	"os"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// Set required environment variable
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error = %v", err)
	}

	if cfg.OpenAIAPIKey != "test-api-key" {
		t.Errorf("OpenAIAPIKey = %v, want %v", cfg.OpenAIAPIKey, "test-api-key")
	}

	if cfg.Model != "whisper-1" {
		t.Errorf("Model = %v, want %v", cfg.Model, "whisper-1")
	}

	if cfg.Language != "en" {
		t.Errorf("Language = %v, want %v", cfg.Language, "en")
	}
}

func TestLoad_CustomValues(t *testing.T) {
	// Set all environment variables
	os.Setenv("OPENAI_API_KEY", "custom-key")
	os.Setenv("STT_MODEL", "custom-model")
	os.Setenv("STT_LANGUAGE", "es")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("STT_MODEL")
		os.Unsetenv("STT_LANGUAGE")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error = %v", err)
	}

	if cfg.OpenAIAPIKey != "custom-key" {
		t.Errorf("OpenAIAPIKey = %v, want %v", cfg.OpenAIAPIKey, "custom-key")
	}

	if cfg.Model != "custom-model" {
		t.Errorf("Model = %v, want %v", cfg.Model, "custom-model")
	}

	if cfg.Language != "es" {
		t.Errorf("Language = %v, want %v", cfg.Language, "es")
	}
}

func TestLoad_MissingAPIKey(t *testing.T) {
	// Ensure API key is not set
	os.Unsetenv("OPENAI_API_KEY")

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error for missing API key, got nil")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "use default when env not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
		{
			name:         "use env value when set",
			key:          "TEST_VAR_SET",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnvOrDefault(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
