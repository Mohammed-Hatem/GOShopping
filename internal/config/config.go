package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
)

func NewConfig() error {
	// Load environment variables from .env file
	envPaths := []string{
		".env",                      // Current directory
		"../.env",                   // One level up
		"../../.env",                // Two levels up (from cmd/api/)
		filepath.Join("..", ".env"), // Alternative path format
	}

	var envLoaded bool
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		return fmt.Errorf("Canoot load the .env file")
	}

	return nil
}
