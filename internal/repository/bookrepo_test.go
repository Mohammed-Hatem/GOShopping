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

	books, err := repo.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}

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

	allBooks, err := repo.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testIsbn := allBooks[0].Isbn
	book, err := repo.GetBookByIsbn(testIsbn)
	if err != nil {
		t.Fatal(err)
	}

	if book == nil {
		t.Fatal("expected a book, got nil")
	}

	if book.Isbn != testIsbn {
		t.Errorf("expected ISBN %s, got %s", testIsbn, book.Isbn)
	}
	fmt.Println(*book)

	nonExistentIsbn := "9999999999999"
	emptyBook, err := repo.GetBookByIsbn(nonExistentIsbn)
	if err != nil {
		t.Fatal(err)
	}

	if emptyBook != nil {
		t.Errorf("expected nil for non-existent ISBN, got %+v", emptyBook)
	}
}

func TestGetBookByTitle(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Test with a valid title
	allBooks, err := repo.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testTitle := allBooks[0].Title
	books, err := repo.GetBookByTitle(testTitle)
	if err != nil {
		t.Fatal(err)
	}

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
	emptyBooks, err := repo.GetBookByTitle(nonExistentTitle)
	if err != nil {
		t.Fatal(err)
	}

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

	allBooks, err := repo.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testPubYear := allBooks[0].PubYear
	books, err := repo.GetBookByPubyear(testPubYear)
	if err != nil {
		t.Fatal(err)
	}

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
	emptyBooks, err := repo.GetBookByPubyear(nonExistentYear)
	if err != nil {
		t.Fatal(err)
	}

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
	allBooks, err := repo.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	testCategory := allBooks[0].Category
	books, err := repo.GetBookByCategory(testCategory)
	if err != nil {
		t.Fatal(err)
	}

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
	emptyBooks, err := repo.GetBookByCategory(nonExistentCategory)
	if err != nil {
		t.Fatal(err)
	}

	if emptyBooks == nil {
		t.Error("should return empty slice, not nil")
	}

	if len(emptyBooks) != 0 {
		t.Errorf("expected empty slice for non-existent category, got %d books", len(emptyBooks))
	}
}

func TestGetBookByAuthor(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	allBooks, err := repo.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(allBooks) == 0 {
		t.Skip("No books in database to test with")
		return
	}

	// Test with a valid author name from seed data
	testAuthor := "Robert C. Martin"
	books, err := repo.GetBookByAuthor(testAuthor)
	if err != nil {
		t.Fatal(err)
	}

	if books == nil {
		t.Error("books should not be nil")
	}

	if len(books) == 0 {
		t.Error("expected at least one book with the given author")
	}

	for _, book := range books {
		fmt.Println(book)
	}

	// Test with a non-existent author
	nonExistentAuthor := "Non-Existent Author Name 12345"
	emptyBooks, err := repo.GetBookByAuthor(nonExistentAuthor)
	if err != nil {
		t.Fatal(err)
	}

	if emptyBooks == nil {
		t.Error("should return empty slice, not nil")
	}

	if len(emptyBooks) != 0 {
		t.Errorf("expected empty slice for non-existent author, got %d books", len(emptyBooks))
	}
}

func TestGetBooksToRestock(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	books, err := repo.GetBooksToRestock()
	if err != nil {
		t.Fatal(err)
	}

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
