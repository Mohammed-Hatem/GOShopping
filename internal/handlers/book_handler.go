package handlers

import (
	"bookstore-project/internal/repository"
	"log"
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
		log.Printf("DEBUG: Failed to decode title: %v", err)
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
