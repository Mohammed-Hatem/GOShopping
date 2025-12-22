package handlers

import (
	"bookstore-project/internal/repository"

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
	books, err := h.repo.GetBookByAuthor(c.Params("name"))

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
