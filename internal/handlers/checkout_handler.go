package handlers

import (
	"errors"
	"strings"
	"time"

	"bookstore-project/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type CheckoutHandler struct {
	cartRepo  *repository.CartRepo
	salesRepo *repository.SalesRepo
	bookRepo  *repository.BookRepo
}

func NewCheckoutHandler(cartRepo *repository.CartRepo, salesRepo *repository.SalesRepo, bookRepo *repository.BookRepo) *CheckoutHandler {
	return &CheckoutHandler{
		cartRepo:  cartRepo,
		salesRepo: salesRepo,
		bookRepo:  bookRepo,
	}
}

type checkoutRequest struct {
	CreditCardNo string `json:"credit_card_no"`
	ExpiryDate   string `json:"expiry_date"` // Format: YYYY-MM-DD
}

// validateCreditCard validates credit card number format (basic validation)
func validateCreditCard(cardNo string) error {
	cardNo = strings.ReplaceAll(cardNo, " ", "")
	cardNo = strings.ReplaceAll(cardNo, "-", "")

	if len(cardNo) < 13 || len(cardNo) > 19 {
		return errors.New("invalid credit card number length")
	}

	// Basic Luhn algorithm check would go here, but for simplicity we'll just check length and digits
	for _, char := range cardNo {
		if char < '0' || char > '9' {
			return errors.New("credit card number must contain only digits")
		}
	}

	return nil
}

// Checkout processes the checkout and creates an order
func (h *CheckoutHandler) Checkout(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req checkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate credit card
	if err := validateCreditCard(req.CreditCardNo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Parse expiry date
	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid expiry date format, expected YYYY-MM-DD",
		})
	}

	// Check if expiry date is in the past
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	expiryStart := time.Date(expiryDate.Year(), expiryDate.Month(), expiryDate.Day(), 0, 0, 0, 0, expiryDate.Location())
	if expiryStart.Before(todayStart) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "credit card has expired",
		})
	}

	// Get cart
	cart, err := h.cartRepo.GetCartByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart",
		})
	}

	// Get cart items
	cartItems, err := h.cartRepo.GetCartItems(cart.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart items",
		})
	}

	if len(cartItems) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cart is empty",
		})
	}

	// Validate stock availability and calculate total
	totalAmount := 0.0
	for _, item := range cartItems {
		// Check stock availability
		if item.StockQuantity < item.Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "insufficient stock for book: " + item.ISBN,
			})
		}
		totalAmount += item.SellingPrice * float64(item.Quantity)
	}

	// Create sales order
	orderID, err := h.salesRepo.CreateSalesOrder(username, totalAmount, req.CreditCardNo, expiryDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create order: " + err.Error(),
		})
	}

	// Create order items and update stock
	for _, item := range cartItems {
		// Create order item
		err = h.salesRepo.CreateOrderItem(orderID, item.ISBN, item.Quantity, item.SellingPrice)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create order item for " + item.ISBN + ": " + err.Error(),
			})
		}

		// Update book stock (deduct purchased quantity)
		book, err := h.bookRepo.GetBookByIsbn(item.ISBN)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to get book",
			})
		}
		if book == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "book not found: " + item.ISBN,
			})
		}

		newStock := book.StockQuantity - item.Quantity
		err = h.bookRepo.UpdateBookStock(item.ISBN, newStock)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update stock",
			})
		}
	}

	// Clear cart after successful checkout
	err = h.cartRepo.ClearCart(cart.ID)
	if err != nil {
		// note to self: add retry if there is enough time
	}

	// Get the created order
	order, err := h.salesRepo.GetSalesOrder(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch order",
		})
	}
	if order == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "order was created but could not be retrieved",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "order created successfully",
		"order":    order,
		"order_id": orderID,
	})
}
