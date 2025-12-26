package handlers

import (
	"strings"

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

