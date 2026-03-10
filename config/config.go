package config

import (
	"errors"
	"os"
)

// AppConfig holds validated application-level configuration.
type AppConfig struct {
	JWTSecretKey string
}

// LoadAppConfig reads and validates required environment variables.
// It must be called after the .env file has been loaded.
func LoadAppConfig() (*AppConfig, error) {
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		return nil, errors.New("JWT_SECRET_KEY environment variable is not set")
	}
	return &AppConfig{JWTSecretKey: jwtKey}, nil
}
