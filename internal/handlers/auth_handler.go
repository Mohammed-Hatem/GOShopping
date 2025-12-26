package handlers

import (
	"bookstore-project/internal/middleware"
	"bookstore-project/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	customerRepo *repository.CustomerRepo
	adminRepo    *repository.AdminRepo
}

func NewAuthHandler(customerRepo *repository.CustomerRepo, adminRepo *repository.AdminRepo) *AuthHandler {
	return &AuthHandler{
		customerRepo: customerRepo,
		adminRepo:    adminRepo,
	}
}

type authLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // "customer" or "admin"
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req authLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Role == "admin" {
		// Admin login
		admin, err := h.adminRepo.GetAdminByUsername(req.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch admin",
			})
		}
		if admin == nil || admin.Password != req.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid username or password",
			})
		}

		token, err := middleware.GenerateAuthToken(req.Username, "admin")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate token",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"username": req.Username,
			"role":     "admin",
			"token":    token,
		})
	} else {
		customer, err := h.customerRepo.GetCustomerByUsername(req.Username)
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
			"role":     "customer",
			"token":    token,
		})
	}
}
