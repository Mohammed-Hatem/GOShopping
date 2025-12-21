package repository

import (
	"bookstore-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type BookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (B *BookRepo) GetAllBooks() []models.Book {
	var books []models.Book
	err := B.db.Select(&books, "SELECT * FROM book")
	if err != nil {
		return nil
	}

	return books
}

func (B *BookRepo) GetBookByIsbn(isbn string) []models.Book {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE isbn = $1", isbn)

	if err != nil {
		return []models.Book{}
	}

	return books
}

func (B *BookRepo) GetBookByTitle(Title string) []models.Book {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE title = $1", Title)

	if err != nil {
		return []models.Book{}
	}

	return books
}

func (B *BookRepo) GetBookByPubyear(PubYear int) []models.Book {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE publication_year = $1", PubYear)

	if err != nil {
		return []models.Book{}
	}

	return books
}

func (B *BookRepo) GetBookByCategory(category string) []models.Book {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE category = $1", category)

	if err != nil {
		return []models.Book{}
	}

	return books
}

func (B *BookRepo) GetBookByAuthor(author string) []models.Book {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE author = $1", author)

	if err != nil {
		return []models.Book{}
	}

	return books
}

func (B *BookRepo) GetBooksToRestock() []models.Book {
	books := []models.Book{}

	err := B.db.Select(&books, "SELECT * FROM book WHERE stock_quantity <= threshold")

	if err != nil {
		return []models.Book{}
	}

	return books
}
