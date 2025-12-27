package database

import (
	"embed"
	"fmt"
	"path"
	"sort"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// createMigrationsTable creates a table to track which migrations have been run
func createMigrationsTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

// isMigrationApplied checks if a migration has already been run
func isMigrationApplied(db *sqlx.DB, version string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// markMigrationApplied records that a migration has been run
func markMigrationApplied(db *sqlx.DB, version string) error {
	query := "INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT (version) DO NOTHING"
	_, err := db.Exec(query, version)
	return err
}

func RunMigrations(db *sqlx.DB) error {
	// Create migrations tracking table
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		// Check if migration has already been applied
		applied, err := isMigrationApplied(db, name)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", name, err)
		}
		if applied {
			continue // Skip already applied migrations
		}

		// Run the migration
		filePath := path.Join("migrations", name)
		sqlBytes, err := migrationsFS.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("exec migration %s: %w", name, err)
		}

		// Mark migration as applied
		if err := markMigrationApplied(db, name); err != nil {
			return fmt.Errorf("failed to mark migration %s as applied: %w", name, err)
		}
	}

	return nil
}
