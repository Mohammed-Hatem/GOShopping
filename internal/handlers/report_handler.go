package handlers

import (
	"bookstore-project/internal/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	repo *repository.ReportRepo
}

func NewReportHandler(repo *repository.ReportRepo) *ReportHandler {
	return &ReportHandler{repo: repo}
}

// GetMonthlySales gets total sales for the previous month
func (h *ReportHandler) GetMonthlySales(c *fiber.Ctx) error {
	result, err := h.repo.GetMonthlySales()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch monthly sales",
		})
	}
	return c.JSON(result)
}

// GetDailySales gets total sales for a specific date
func (h *ReportHandler) GetDailySales(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "date parameter is required (format: YYYY-MM-DD)",
		})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid date format, expected YYYY-MM-DD",
		})
	}

	result, err := h.repo.GetDailySales(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch daily sales",
		})
	}
	return c.JSON(result)
}

// GetTopCustomers gets top 5 customers by total purchase for the last 3 months
func (h *ReportHandler) GetTopCustomers(c *fiber.Ctx) error {
	months := 3 // Default to 3 months as specified in requirements
	monthsStr := c.Query("months")
	if monthsStr != "" {
		parsedMonths, err := strconv.Atoi(monthsStr)
		if err == nil && parsedMonths > 0 {
			months = parsedMonths
		}
	}

	customers, err := h.repo.GetTopCustomers(months)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch top customers",
		})
	}
	return c.JSON(customers)
}

// GetTopSellingBooks gets top 10 selling books by quantity for the last 3 months
func (h *ReportHandler) GetTopSellingBooks(c *fiber.Ctx) error {
	months := 3 // Default to 3 months as specified in requirements
	monthsStr := c.Query("months")
	if monthsStr != "" {
		parsedMonths, err := strconv.Atoi(monthsStr)
		if err == nil && parsedMonths > 0 {
			months = parsedMonths
		}
	}

	books, err := h.repo.GetTopSellingBooks(months)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch top selling books",
		})
	}
	return c.JSON(books)
}

// GetBookOrderCount gets the total number of times a specific book has been ordered
func (h *ReportHandler) GetBookOrderCount(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn parameter is required",
		})
	}

	result, err := h.repo.GetBookOrderCount(isbn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch book order count",
		})
	}
	return c.JSON(result)
}
