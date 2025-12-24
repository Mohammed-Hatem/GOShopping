package repository

import (
	"bookstore-project/internal/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type BookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (B *BookRepo) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := B.db.Select(&books, "SELECT * FROM book")
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepo) GetBookByIsbn(isbn string) (*models.Book, error) {
	var book models.Book
	err := r.db.Get(&book, "SELECT * FROM book WHERE isbn = $1", isbn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (B *BookRepo) GetBookByTitle(title string) ([]models.Book, error) {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE LOWER(title) = LOWER($1)", title)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (B *BookRepo) GetBookByPubyear(pubYear int) ([]models.Book, error) {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE publication_year = $1", pubYear)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (B *BookRepo) GetBookByCategory(category string) ([]models.Book, error) {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE LOWER(category) = LOWER($1)", category)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (B *BookRepo) GetBookByAuthor(author string) ([]models.Book, error) {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE LOWER(author_name) = LOWER($1)", author)

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (B *BookRepo) GetBooksToRestock() ([]models.Book, error) {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE stock_quantity <= threshold")

	if err != nil {
		return nil, err
	}

	return books, nil
}
