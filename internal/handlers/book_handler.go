package handlers

import (
	"bookstore-project/internal/models"
	"bookstore-project/internal/repository"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	repo *repository.BookRepo
}

func NewBookHandler(repo *repository.BookRepo) *BookHandler {
	return &BookHandler{repo: repo}
}

func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	books, err := h.repo.GetAllBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch books",
		})
	}
	return c.JSON(books)
}

// Search books by publisher
func (h *BookHandler) GetBookByPublisher(c *fiber.Ctx) error {
	publisherId, err := strconv.Atoi(c.Params("publisher_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid publisher ID",
		})
	}

	books, err := h.repo.GetBookByPublisher(publisherId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch books",
		})
	}

	if len(books) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no books found for this publisher",
		})
	}

	return c.JSON(books)
}

// Admin only: Add new book
func (h *BookHandler) AddBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate required fields
	if book.Isbn == "" || book.Title == "" || book.AuthorName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn, title, and author_name are required",
		})
	}

	// Validate category
	validCategories := []string{"Science", "Art", "Religion", "History", "Geography"}
	isValidCategory := false
	for _, cat := range validCategories {
		if book.Category == cat {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "category must be one of: Science, Art, Religion, History, Geography",
		})
	}

	// Validate threshold
	if book.Threshold <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "threshold must be positive",
		})
	}

	// Validate stock quantity
	if book.StockQuantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "stock_quantity cannot be negative",
		})
	}

	err := h.repo.AddBook(book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add book",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "book added successfully",
		"book":    book,
	})
}

// Admin only: Update existing book
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn is required",
		})
	}

	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Set ISBN from URL parameter
	book.Isbn = isbn

	// Validate category if provided
	if book.Category != "" {
		validCategories := []string{"Science", "Art", "Religion", "History", "Geography"}
		isValidCategory := false
		for _, cat := range validCategories {
			if book.Category == cat {
				isValidCategory = true
				break
			}
		}
		if !isValidCategory {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "category must be one of: Science, Art, Religion, History, Geography",
			})
		}
	}

	// Validate threshold if provided
	if book.Threshold <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "threshold must be positive",
		})
	}

	err := h.repo.UpdateBook(book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update book",
		})
	}

	return c.JSON(fiber.Map{
		"message": "book updated successfully",
		"book":    book,
	})
}

// Admin only: Update book stock quantity
func (h *BookHandler) UpdateBookStock(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn is required",
		})
	}

	type StockUpdate struct {
		Quantity int `json:"quantity"`
	}

	var update StockUpdate
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if update.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "stock quantity cannot be negative",
		})
	}

	err := h.repo.UpdateBookStock(isbn, update.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update stock quantity",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "stock quantity updated successfully",
		"isbn":     isbn,
		"quantity": update.Quantity,
	})
}
func (h *BookHandler) GetBookByIsbn(c *fiber.Ctx) error {
	book, err := h.repo.GetBookByIsbn(c.Params("isbn"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch book",
		})
	}

	if book == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "book not found",
		})
	}

	return c.JSON(book)
}

func (h *BookHandler) GetBookByAuthor(c *fiber.Ctx) error {

	author := c.Params("name")
	decodedAuthor, err := url.QueryUnescape(author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid author encoding",
		})
	}
	books, err := h.repo.GetBookByAuthor(decodedAuthor)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch book",
		})
	}

	if len(books) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no books found",
		})
	}

	return c.JSON(books)

}

func (h *BookHandler) GetBookByPubyear(c *fiber.Ctx) error {
	pubyear, err := strconv.Atoi(c.Params("year"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid publication year",
		})
	}

	books, err := h.repo.GetBookByPubyear(pubyear)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch book",
		})
	}

	if len(books) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no books found",
		})
	}

	return c.JSON(books)

}

func (h *BookHandler) GetBookByTitle(c *fiber.Ctx) error {
	title := c.Params("title")

	//  DECODE title TO RESOLVE SPACES AND SPECIAL CHARACTERS
	decodedTitle, err := url.QueryUnescape(title)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid title encoding",
		})
	}

	if decodedTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid title",
		})
	}

	books, err := h.repo.GetBookByTitle(decodedTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch book",
		})
	}

	if len(books) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no books found",
		})
	}

	return c.JSON(books)

}

func (h *BookHandler) GetBookByCategory(c *fiber.Ctx) error {
	category := c.Params("category")

	decodedCategory, err := url.QueryUnescape(category)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid category encoding",
		})
	}

	if decodedCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid category",
		})
	}

	books, err := h.repo.GetBookByCategory(decodedCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch books",
		})
	}

	if len(books) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no books found",
		})
	}

	return c.JSON(books)

}
