package repository

import (
	"bookstore-project/internal/database"
	"fmt"
	"testing"

	"bookstore-project/internal/config"
)

func setupTestDB(t *testing.T) (*BookRepo, func()) {
	err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewBookRepo(db)
	cleanup := func() {
		db.Close()
	}

	return repo, cleanup
}

func TestGetAllBooks(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	books := repo.GetAllBooks()

	if books == nil {
		t.Error("books should not be nil")
	}

	if len(books) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	for _, book := range books {
		fmt.Println(book)
	}
}

func TestGetBookByIsbn(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	allBooks := repo.GetAllBooks()
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testIsbn := allBooks[0].Isbn
	books := repo.GetBookByIsbn(testIsbn)

	if books == nil {
		t.Error("books should not be nil")
	}

	if len(books) == 0 {
		t.Error("expected at least one book with the given ISBN")
	}

	for _, book := range books {
		if book.Isbn != testIsbn {
			t.Errorf("expected ISBN %s, got %s", testIsbn, book.Isbn)
		}
		fmt.Println(book)
	}

	nonExistentIsbn := "9999999999999"
	emptyBooks := repo.GetBookByIsbn(nonExistentIsbn)

	if emptyBooks == nil {
		t.Error("should return empty slice, not nil")
	}

	if len(emptyBooks) != 0 {
		t.Errorf("expected empty slice for non-existent ISBN, got %d books", len(emptyBooks))
	}
}

func TestGetBookByTitle(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Test with a valid title
	allBooks := repo.GetAllBooks()
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testTitle := allBooks[0].Title
	books := repo.GetBookByTitle(testTitle)

	if books == nil {
		t.Error("books should not be nil")
	}

	if len(books) == 0 {
		t.Error("expected at least one book with the given title")
	}

	for _, book := range books {
		if book.Title != testTitle {
			t.Errorf("expected title %s, got %s", testTitle, book.Title)
		}
		fmt.Println(book)
	}

	nonExistentTitle := "Non-Existent Book Title 12345"
	emptyBooks := repo.GetBookByTitle(nonExistentTitle)

	if emptyBooks == nil {
		t.Error("should return empty slice, not nil")
	}

	if len(emptyBooks) != 0 {
		t.Errorf("expected empty slice for non-existent title, got %d books", len(emptyBooks))
	}
}

func TestGetBookByPubyear(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	allBooks := repo.GetAllBooks()
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testPubYear := allBooks[0].PubYear
	books := repo.GetBookByPubyear(testPubYear)

	if books == nil {
		t.Error("books should not be nil")
	}

	for _, book := range books {
		if book.PubYear != testPubYear {
			t.Errorf("expected publication year %d, got %d", testPubYear, book.PubYear)
		}
		fmt.Println(book)
	}

	nonExistentYear := 9999
	emptyBooks := repo.GetBookByPubyear(nonExistentYear)

	if emptyBooks == nil {
		t.Error("should return empty slice, not nil")
	}

	if len(emptyBooks) != 0 {
		t.Errorf("expected empty slice for non-existent publication year, got %d books", len(emptyBooks))
	}
}

func TestGetBookByCategory(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Test with a valid category
	allBooks := repo.GetAllBooks()
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testCategory := allBooks[0].Category
	books := repo.GetBookByCategory(testCategory)

	if books == nil {
		t.Error("books should not be nil")
	}

	for _, book := range books {
		if book.Category != testCategory {
			t.Errorf("expected category %s, got %s", testCategory, book.Category)
		}
		fmt.Println(book)
	}

	nonExistentCategory := "NonExistentCategory123"
	emptyBooks := repo.GetBookByCategory(nonExistentCategory)

	if emptyBooks == nil {
		t.Error("should return empty slice, not nil")
	}

	if len(emptyBooks) != 0 {
		t.Errorf("expected empty slice for non-existent category, got %d books", len(emptyBooks))
	}
}

func TestGetBooksToRestock(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	books := repo.GetBooksToRestock()

	if books == nil {
		t.Error("books should not be nil")
	}

	// Verify all returned books have stock_quantity <= threshold
	for _, book := range books {
		if book.StockQuantity > book.Threshold {
			t.Errorf("book %s has stock_quantity %d which is greater than threshold %d",
				book.Isbn, book.StockQuantity, book.Threshold)
		}
		fmt.Printf("Book: %s, Stock: %d, Threshold: %d\n", book.Title, book.StockQuantity, book.Threshold)
	}
}
