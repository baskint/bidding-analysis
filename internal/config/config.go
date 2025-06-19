package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	OpenAI   OpenAIConfig
	JWT      JWTConfig
	GCP      GCPConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port     int
	GRPCPort int
}

// OpenAIConfig holds OpenAI API settings
type OpenAIConfig struct {
	APIKey string
}

// JWTConfig holds JWT settings
type JWTConfig struct {
	Secret string
}

// GCPConfig holds Google Cloud Platform settings
type GCPConfig struct {
	ProjectID string
	Region    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		// Don't error if .env doesn't exist (production might use env vars directly)
		fmt.Println("No .env file found, using environment variables")
	}

	cfg := &Config{}

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.Name = getEnv("DB_NAME", "bidding_analysis")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

	// Server configuration
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	cfg.Server.GRPCPort = getEnvAsInt("GRPC_PORT", 9090)

	// OpenAI configuration
	cfg.OpenAI.APIKey = getEnv("OPENAI_API_KEY", "")

	// JWT configuration
	cfg.JWT.Secret = getEnv("JWT_SECRET", "default-secret-change-in-production")

	// GCP configuration
	cfg.GCP.ProjectID = getEnv("GOOGLE_CLOUD_PROJECT", "")
	cfg.GCP.Region = getEnv("GCP_REGION", "us-central1")

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that required configuration is present
func (c *Config) validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if c.OpenAI.APIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}

	return nil
}

// DatabaseURL returns the PostgreSQL connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode)
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
