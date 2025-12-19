package main

import (
	"log"
	"path/filepath"

	"bookstore-project/internal/database"

	"github.com/joho/godotenv"
)

func main() {
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
		log.Printf("Warning: .env file not found in common locations, using system environment variables")
	}

	// cnnect to database and verify connection
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")
}
