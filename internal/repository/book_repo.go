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

// AddBook adds a new book to the database
func (r *BookRepo) AddBook(book models.Book) error {
	query := `
		INSERT INTO book (isbn, title, publication_year, selling_price, category, author_name, stock_quantity, threshold, publisher_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, book.Isbn, book.Title, book.PubYear, book.SellingPrice,
		book.Category, book.AuthorName, book.StockQuantity, book.Threshold, book.PublisherId)
	return err
}

// UpdateBook updates an existing book
func (r *BookRepo) UpdateBook(book models.Book) error {
	query := `
		UPDATE book 
		SET title = $2, publication_year = $3, selling_price = $4, category = $5, 
			author_name = $6, threshold = $7, publisher_id = $8
		WHERE isbn = $1`

	_, err := r.db.Exec(query, book.Isbn, book.Title, book.PubYear, book.SellingPrice,
		book.Category, book.AuthorName, book.Threshold, book.PublisherId)
	return err
}

// UpdateBookStock updates the stock quantity of a book
func (r *BookRepo) UpdateBookStock(isbn string, quantity int) error {
	query := `UPDATE book SET stock_quantity = $2 WHERE isbn = $1`
	_, err := r.db.Exec(query, isbn, quantity)
	return err
}

// GetBookByPublisher retrieves books by publisher ID
func (r *BookRepo) GetBookByPublisher(publisherId int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Select(&books, "SELECT * FROM book WHERE publisher_id = $1", publisherId)
	if err != nil {
		return nil, err
	}
	return books, nil
}
