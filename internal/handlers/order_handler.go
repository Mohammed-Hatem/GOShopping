package handlers

import (
	"bookstore-project/internal/models"
	"bookstore-project/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	repo *repository.OrderRepo
}

func NewOrderHandler(repo *repository.OrderRepo) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// Admin only: Get all publisher orders
func (h *OrderHandler) GetAllPublisherOrders(c *fiber.Ctx) error {
	orders, err := h.repo.GetAllPublisherOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch publisher orders",
		})
	}
	return c.JSON(orders)
}

// Admin only: Get publisher order by ID
func (h *OrderHandler) GetPublisherOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	order, err := h.repo.GetPublisherOrder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch publisher order",
		})
	}

	if order == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "publisher order not found",
		})
	}

	return c.JSON(order)
}

// Admin only: Place publisher order manually
func (h *OrderHandler) PlacePublisherOrder(c *fiber.Ctx) error {
	var order models.PublisherOrder
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate required fields
	if order.ISBN == "" || order.Quantity <= 0 || order.AdminID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isbn, positive quantity, and admin_id are required",
		})
	}

	err := h.repo.PlacePublisherOrder(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to place publisher order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "publisher order placed successfully",
		"order":   order,
	})
}

// Admin only: Confirm publisher order
func (h *OrderHandler) ConfirmPublisherOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	err = h.repo.ConfirmPublisherOrder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to confirm publisher order",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "publisher order confirmed successfully",
		"order_id": id,
	})
}

// Admin only: Get pending publisher orders
func (h *OrderHandler) GetPendingPublisherOrders(c *fiber.Ctx) error {
	orders, err := h.repo.GetPendingPublisherOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch pending publisher orders",
		})
	}
	return c.JSON(orders)
}

// Admin only: Update publisher order status
func (h *OrderHandler) UpdatePublisherOrderStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if request.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status is required",
		})
	}

	err = h.repo.UpdatePublisherOrderStatus(id, request.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Fetch updated order to return
	order, err := h.repo.GetPublisherOrder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch updated order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "order status updated successfully",
		"order":   order,
	})
}