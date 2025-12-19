package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New() (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	return sqlx.Connect("postgres", dsn)
}

func ExecSQLFile(db *sqlx.DB, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(data))
	return err
}
