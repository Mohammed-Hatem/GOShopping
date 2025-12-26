package handlers

import (
	"strings"
	"time"

	"bookstore-project/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartRepo  *repository.CartRepo
	orderRepo *repository.OrderRepo
}

func NewCartHandler(cartRepo *repository.CartRepo, orderRepo *repository.OrderRepo) *CartHandler {
	return &CartHandler{cartRepo: cartRepo, orderRepo: orderRepo}
}

type addToCartRequest struct {
	ISBN     string `json:"isbn"`
	Quantity int    `json:"quantity"`
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req addToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.ISBN = strings.TrimSpace(req.ISBN)
	if req.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "isbn is required"})
	}
	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "quantity must be positive"})
	}

	err := h.cartRepo.AddToCart(username, req.ISBN, req.Quantity)
	if err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add to cart"})
	}

	summary, err := h.cartRepo.GetCartSummary(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load cart"})
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

func (h *CartHandler) ViewCart(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	summary, err := h.cartRepo.GetCartSummary(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load cart"})
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	isbn := strings.TrimSpace(c.Params("isbn"))
	if isbn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "isbn is required"})
	}

	if err := h.cartRepo.RemoveFromCart(username, isbn); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove item"})
	}

	summary, err := h.cartRepo.GetCartSummary(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load cart"})
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

type checkoutRequest struct {
	CreditCardNo string `json:"credit_card_no"`
	ExpiryDate   string `json:"expiry_date"` // YYYY-MM-DD
}

func (h *CartHandler) Checkout(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req checkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	req.CreditCardNo = strings.TrimSpace(req.CreditCardNo)
	req.ExpiryDate = strings.TrimSpace(req.ExpiryDate)

	if req.CreditCardNo == "" || req.ExpiryDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit_card_no and expiry_date are required"})
	}
	if !isValidCreditCardNumber(req.CreditCardNo) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid credit card number"})
	}

	expiry, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "expiry_date must be YYYY-MM-DD"})
	}
	// Treat expiry as end-of-day local time to avoid off-by-one.
	expiry = time.Date(expiry.Year(), expiry.Month(), expiry.Day(), 23, 59, 59, 0, time.Local)
	if time.Now().After(expiry) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit card expired"})
	}

	res, err := h.orderRepo.CheckoutCart(username, req.CreditCardNo, expiry)
	if err != nil {
		status := fiber.StatusBadRequest
		if isInternalCheckoutError(err) {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func isInternalCheckoutError(err error) bool {
	if err == nil {
		return false
	}
	// Business-rule errors should be 400; everything else bubbles as 500.
	s := err.Error()
	switch {
	case strings.Contains(s, "cart is empty"):
		return false
	case strings.Contains(s, "insufficient stock"):
		return false
	case strings.Contains(s, "book not found"):
		return false
	case strings.Contains(s, "required"):
		return false
	default:
		return true
	}
}

func isValidCreditCardNumber(card string) bool {
	// Accept digits only.
	digits := make([]int, 0, len(card))
	for _, r := range card {
		if r < '0' || r > '9' {
			return false
		}
		digits = append(digits, int(r-'0'))
	}
	if len(digits) < 13 || len(digits) > 19 {
		return false
	}

	// Luhn check.
	sum := 0
	double := false
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}

