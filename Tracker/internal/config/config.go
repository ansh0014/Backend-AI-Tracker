package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// ConfigError represents a configuration-specific error
type ConfigError struct {
	Field   string
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("config error in %s: %s: %v", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("config error in %s: %s", e.Field, e.Message)
}

// Config holds the application configuration
type Config struct {
	// Server
	Port string
	Env  string

	// MongoDB
	MongoURI        string
	MongoDBName     string
	MongoCollection string

	// Gemini AI
	GeminiApiKey string
	GeminiModel  string

	// Logging
	LogLevel string
}

// LoadConfig loads the configuration from .env file
func LoadConfig() (*Config, error) {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return nil, &ConfigError{
			Field:   "working_directory",
			Message: "failed to get current directory",
			Err:     err,
		}
	}

	// Look for .env file in the current directory and parent directories
	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Try parent directory
		envPath = filepath.Join(dir, "..", ".env")
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			return nil, &ConfigError{
				Field:   "env_file",
				Message: ".env file not found in current or parent directory",
			}
		}
	}

	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		return nil, &ConfigError{
			Field:   "env_file",
			Message: "failed to load .env file",
			Err:     err,
		}
	}

	config := &Config{
		// Server
		Port: getEnvOrDefault("PORT", "8080"),
		Env:  getEnvOrDefault("ENV", "development"),

		// MongoDB
		MongoURI:        getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName:     getEnvOrDefault("MONGODB_DB", "activity_tracker"),
		MongoCollection: getEnvOrDefault("MONGODB_COLLECTION", "activities"),

		// Gemini AI
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
		GeminiModel:  getEnvOrDefault("GEMINI_MODEL", "gemini-pro"),

		// Logging
		LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	var errors []string

	// Validate Server configuration
	if c.Port == "" {
		errors = append(errors, "PORT is required")
	}
	if c.Env == "" {
		errors = append(errors, "ENV is required")
	}
	if !isValidEnvironment(c.Env) {
		errors = append(errors, fmt.Sprintf("invalid ENV value: %s (must be one of: development, production, testing)", c.Env))
	}

	// Validate MongoDB configuration
	if c.MongoURI == "" {
		errors = append(errors, "MONGODB_URI is required")
	}
	if c.MongoDBName == "" {
		errors = append(errors, "MONGODB_DB is required")
	}
	if c.MongoCollection == "" {
		errors = append(errors, "MONGODB_COLLECTION is required")
	}

	// Validate Gemini AI configuration
	if c.GeminiApiKey == "" {
		errors = append(errors, "GEMINI_API_KEY is required")
	}
	if c.GeminiModel == "" {
		errors = append(errors, "GEMINI_MODEL is required")
	}

	// Validate Logging configuration
	if !isValidLogLevel(c.LogLevel) {
		errors = append(errors, fmt.Sprintf("invalid LOG_LEVEL: %s (must be one of: debug, info, warn, error)", c.LogLevel))
	}

	if len(errors) > 0 {
		return &ConfigError{
			Field:   "validation",
			Message: fmt.Sprintf("configuration validation failed: %s", strings.Join(errors, "; ")),
		}
	}

	return nil
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isValidEnvironment checks if the environment value is valid
func isValidEnvironment(env string) bool {
	validEnvs := map[string]bool{
		"development": true,
		"production":  true,
		"testing":     true,
	}
	return validEnvs[env]
}

// isValidLogLevel checks if the log level is valid
func isValidLogLevel(level string) bool {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	return validLevels[level]
}

// GetMongoURI returns the MongoDB URI
func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}

// GetMongoDBName returns the MongoDB database name
func GetMongoDBName() string {
	return os.Getenv("MONGODB_DB")
}

// GetMongoCollectionName returns the MongoDB collection name
func GetMongoCollectionName() string {
	return os.Getenv("MONGODB_COLLECTION")
}

// GetGeminiApiKey returns the Gemini API key
func GetGeminiApiKey() string {
	return os.Getenv("GEMINI_API_KEY")
}