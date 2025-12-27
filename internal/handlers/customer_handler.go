package handlers

import (
	"strconv"
	"strings"

	"bookstore-project/internal/models"
	"bookstore-project/internal/middleware"
	"bookstore-project/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	repo *repository.CustomerRepo
}

func NewCustomerHandler(repo *repository.CustomerRepo) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

type signupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"shipping_address"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type updateProfileRequest struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Address        *string `json:"shipping_address"`
	Password       *string `json:"password"`
	CurrentPassword *string `json:"current_password"`
}

func (h *CustomerHandler) Signup(c *fiber.Ctx) error {
	var req signupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Password == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username, password, and email are required",
		})
	}

	err := h.repo.CreateCustomer(req.Username, req.Password, req.FirstName, req.LastName, req.Email, req.Phone, req.Address)
	if err != nil {
		switch err.Error() {
		case "username already exists", "email already exists":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create customer",
			})
		}
	}

	token, err := middleware.GenerateAuthToken(req.Username, "customer")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"username": req.Username,
		"token":    token,
	})
}

func (h *CustomerHandler) Login(c *fiber.Ctx) error {
	var req loginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "faulty request body",
		})
	}
	customer, err := h.repo.GetCustomerByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch customer",
		})
	}
	if customer == nil || customer.Password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid username or password",
		})
	}

	token, err := middleware.GenerateAuthToken(req.Username, "customer")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": req.Username,
		"token":    token,
	})

}

func (h *CustomerHandler) GetProfile(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	customer, err := h.repo.GetCustomerByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch customer",
		})
	}
	if customer == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "customer not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(customer)
}


func (h *CustomerHandler) UpdateProfile(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")//returns any type (intergace{})

	username, ok := usernameAny.(string) //type assertion
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// If password is being updated, verify current password first
	if req.Password != nil {
		if req.CurrentPassword == nil || *req.CurrentPassword == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "current_password is required when updating password",
			})
		}

		// Get customer to verify current password
		customer, err := h.repo.GetCustomerByUsername(username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch customer",
			})
		}
		if customer == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "customer not found",
			})
		}

		// Verify current password
		if customer.Password != *req.CurrentPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "current password is incorrect",
			})
		}
	}

		err := h.repo.UpdateCustomerProfile(username, req.FirstName, req.LastName, req.Email, req.Phone, req.Address, req.Password)
		if err != nil {
			switch err.Error() {
			case "email already exists":
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": err.Error(),
				})
			case "no fields to update":
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to update profile",
				})
			}
		}

		customer, err := h.repo.GetCustomerByUsername(username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch customer",
			})
		}
		if customer == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "customer not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(customer)
	}

// GetCustomerOrders retrieves all orders for the authenticated customer
func (h *CustomerHandler) GetCustomerOrders(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	orders, err := h.repo.GetOrdersByCustomer(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

// GetOrderDetails retrieves detailed information about a specific order
func (h *CustomerHandler) GetOrderDetails(c *fiber.Ctx) error {
	usernameAny := c.Locals("username")
	username, ok := usernameAny.(string)
	if !ok || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	orderIDParam := c.Params("id")
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	// Get all orders for the user and find the specific one
	orders, err := h.repo.GetOrdersByCustomer(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch orders",
		})
	}

	var order *models.SalesOrder
	for i := range orders {
		if orders[i].OrderID == orderID {
			order = &orders[i]
			break
		}
	}

	if order == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "order not found",
		})
	}

	// Get order items
	items, err := h.repo.GetOrderItemsByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch order items",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"order": order,
		"items": items,
	})
}
