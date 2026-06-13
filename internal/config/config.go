package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	LogLevel string
}

func Load() (*Config, error) {

	// Ignore error if .env doesn't exist.
	// Production environments usually don't use it.
	_ = godotenv.Load()

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost: getEnv("DB_HOST", ""),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", ""),
		DBPass: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", ""),

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {

	if c.DBHost == "" {
		return errors.New("DB_HOST is required")
	}

	if c.DBUser == "" {
		return errors.New("DB_USER is required")
	}

	if c.DBName == "" {
		return errors.New("DB_NAME is required")
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf(
			"invalid LOG_LEVEL '%s' (debug|info|warn|error)",
			c.LogLevel,
		)
	}

	return nil
}

func getEnv(key string, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

// build postgres connection config string
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost,
		c.DBUser,
		c.DBPass,
		c.DBName,
		c.DBPort,
	)
}
