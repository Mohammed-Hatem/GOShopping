package main

import (
	"log"
	"os"
	"strings"

	"bookstore-project/internal/config"
	"bookstore-project/internal/database"
	"bookstore-project/internal/handlers"
	"bookstore-project/internal/middleware"
	"bookstore-project/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	err := config.NewConfig()

	if err != nil {
		panic(err)
	}

	// connect to database and verify connection
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	app := fiber.New()

	// Initialize repositories
	bookRepo := repository.NewBookRepo(db)
	bookHandler := handlers.NewBookHandler(bookRepo)

	customerRepo := repository.NewCustomerRepo(db)
	customerHandler := handlers.NewCustomerHandler(customerRepo)

	adminRepo := repository.NewAdminRepo(db)
	authHandler := handlers.NewAuthHandler(customerRepo, adminRepo)

	orderRepo := repository.NewOrderRepo(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	// Public book search endpoints (available to all users)
	app.Get("/api/books", bookHandler.GetAllBooks)
	app.Get("/api/books/author/:name", bookHandler.GetBookByAuthor)
	app.Get("/api/books/pubyear/:year", bookHandler.GetBookByPubyear)
	app.Get("/api/books/title/:title", bookHandler.GetBookByTitle)
	app.Get("/api/books/category/:category", bookHandler.GetBookByCategory)
	app.Get("/api/books/publisher/:publisher_id", bookHandler.GetBookByPublisher)
	app.Get("/api/books/:isbn", bookHandler.GetBookByIsbn)

	// Authentication endpoints
	app.Post("/api/auth/signup", customerHandler.Signup)
	app.Post("/api/auth/login", authHandler.Login)

	// Customer endpoints
	app.Get("/api/customer/profile", middleware.Protect(), customerHandler.GetProfile)
	app.Patch("/api/customer/profile", middleware.Protect(), customerHandler.UpdateProfile)

	// Admin-only book management endpoints
	app.Post("/api/admin/books", middleware.AdminOnly(), bookHandler.AddBook)
	app.Put("/api/admin/books/:isbn", middleware.AdminOnly(), bookHandler.UpdateBook)
	app.Patch("/api/admin/books/:isbn/stock", middleware.AdminOnly(), bookHandler.UpdateBookStock)

	// Admin-only publisher order endpoints
	app.Get("/api/admin/publisher-orders", middleware.AdminOnly(), orderHandler.GetAllPublisherOrders)
	app.Get("/api/admin/publisher-orders/pending", middleware.AdminOnly(), orderHandler.GetPendingPublisherOrders)
	app.Get("/api/admin/publisher-orders/:id", middleware.AdminOnly(), orderHandler.GetPublisherOrder)
	app.Post("/api/admin/publisher-orders", middleware.AdminOnly(), orderHandler.PlacePublisherOrder)
	app.Put("/api/admin/publisher-orders/:id/confirm", middleware.AdminOnly(), orderHandler.ConfirmPublisherOrder)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3001"
	}
	addr := port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	err = app.Listen(addr)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
