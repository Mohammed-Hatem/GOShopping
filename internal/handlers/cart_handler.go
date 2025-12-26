package handlers

import (
	"bookstore-project/internal/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartRepo *repository.CartRepo
}

func NewCartHandler(cartRepo *repository.CartRepo) *CartHandler {
	return &CartHandler{cartRepo: cartRepo}
}

type addToCartRequest struct {
	ISBN     string `json:"isbn"`
	Quantity int    `json:"quantity"`
}

type updateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

// GetCart retrieves the user's cart with items and total
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get or create cart
	cart, err := h.cartRepo.GetCartByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart",
		})
	}

	// Get cart items with book details
	items, err := h.cartRepo.GetCartItems(cart.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart items",
		})
	}

	// Get cart total
	total, err := h.cartRepo.GetCartTotal(cart.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to calculate cart total",
		})
	}

	return c.JSON(fiber.Map{
		"cart_id": cart.ID,
		"items":   items,
		"total":   total,
	})
}

// AddToCart adds an item to the cart
func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req addToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.ISBN == "" || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn and positive quantity are required",
		})
	}

	// Get or create cart
	cart, err := h.cartRepo.GetCartByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart",
		})
	}

	// Add item to cart
	err = h.cartRepo.AddItemToCart(cart.ID, req.ISBN, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add item to cart",
		})
	}

	// Return updated cart
	return h.GetCart(c)
}

// UpdateCartItem updates the quantity of a cart item
func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	isbn := c.Params("isbn")
	if isbn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn is required",
		})
	}

	var req updateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "quantity must be greater than 0",
		})
	}

	// Get cart
	cart, err := h.cartRepo.GetCartByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart",
		})
	}

	// Update cart item
	err = h.cartRepo.UpdateCartItemQuantity(cart.ID, isbn, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update cart item",
		})
	}

	// Return updated cart
	return h.GetCart(c)
}

// RemoveFromCart removes an item from the cart
func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	isbn := c.Params("isbn")
	if isbn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn is required",
		})
	}

	// Get cart
	cart, err := h.cartRepo.GetCartByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart",
		})
	}

	// Remove item from cart
	err = h.cartRepo.RemoveItemFromCart(cart.ID, isbn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to remove item from cart",
		})
	}

	// Return updated cart
	return h.GetCart(c)
}

// ClearCart clears all items from the cart
func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get cart
	cart, err := h.cartRepo.GetCartByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get cart",
		})
	}

	// Clear cart
	err = h.cartRepo.ClearCart(cart.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to clear cart",
		})
	}

	return c.JSON(fiber.Map{
		"message": "cart cleared successfully",
	})
}
