package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort string

	// Option 1: Supabase / Direct connection string
	DatabaseURL string

	// Option 2: Local PostgreSQL (individual params)
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// LoadConfig reads configuration from environment variables.
// Supports two modes:
//   - DATABASE_URL is set → use it directly (Supabase / cloud PostgreSQL)
//   - DATABASE_URL is empty → build DSN from individual DB_* variables (local PostgreSQL)
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, using system environment variables")
	}

	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "todo_db"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
	}

	if cfg.DatabaseURL != "" {
		log.Info().
			Str("server_port", cfg.ServerPort).
			Str("mode", "supabase/cloud").
			Msg("Configuration loaded (using DATABASE_URL)")
	} else {
		log.Info().
			Str("server_port", cfg.ServerPort).
			Str("mode", "local").
			Str("db_host", cfg.DBHost).
			Str("db_port", cfg.DBPort).
			Str("db_name", cfg.DBName).
			Msg("Configuration loaded (using local PostgreSQL)")
	}

	return cfg
}

// DSN returns the PostgreSQL connection string.
// If DATABASE_URL is set, it returns that directly.
// Otherwise, it builds the DSN from individual parameters.
func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
