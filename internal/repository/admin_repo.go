package repository

import (
	"bookstore-project/internal/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AdminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

// GetAdminByUsername retrieves an admin by username
func (r *AdminRepo) GetAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Get(&admin, "SELECT * FROM administrator WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// CreateAdmin creates a new admin
func (r *AdminRepo) CreateAdmin(username, password string) error {
	query := `INSERT INTO administrator (username, password) VALUES ($1, $2)`
	_, err := r.db.Exec(query, username, password)
	return err
}
