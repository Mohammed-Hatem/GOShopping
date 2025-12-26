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

func RunMigrations(db *sqlx.DB) error {
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
		filePath := path.Join("migrations", name)
		sqlBytes, err := migrationsFS.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("exec migration %s: %w", name, err)
		}
	}

	return nil
}
