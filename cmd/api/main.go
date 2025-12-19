package main

import (
	"log"
	"os"

	"bookstore-project/internal/database"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func execSQLFile(db *sqlx.DB, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(data))
	return err
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables: %v", err)
	}

	// Connect to PostgreSQL using environment variables
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Run DDL (creates tables) from migration file
	if err := execSQLFile(db, "internal/database/migrations/001_create_tables.sql"); err != nil {
		log.Fatalf("failed to run DDL migration: %v", err)
	}
	log.Println("Tables created/ensured successfully")

	// Run seed data inserts
	if err := execSQLFile(db, "internal/database/migrations/003_seed_data.sql"); err != nil {
		log.Fatalf("failed to run seed data: %v", err)
	}
	log.Println("Seed data inserted successfully")
}
