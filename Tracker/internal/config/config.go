package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

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
		return nil, fmt.Errorf("error getting current directory: %v", err)
	}

	// Look for .env file in the current directory and parent directories
	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Try parent directory
		envPath = filepath.Join(dir, "..", ".env")
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			return nil, fmt.Errorf(".env file not found")
		}
	}

	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
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
	if c.GeminiApiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required")
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