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
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Enable CORS for frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Initialize repositories
	bookRepo := repository.NewBookRepo(db)
	bookHandler := handlers.NewBookHandler(bookRepo)

	customerRepo := repository.NewCustomerRepo(db)
	customerHandler := handlers.NewCustomerHandler(customerRepo)

	adminRepo := repository.NewAdminRepo(db)
	authHandler := handlers.NewAuthHandler(customerRepo, adminRepo)

	orderRepo := repository.NewOrderRepo(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	cartRepo := repository.NewCartRepo(db)
	cartHandler := handlers.NewCartHandler(cartRepo)

	salesRepo := repository.NewSalesRepo(db)
	checkoutHandler := handlers.NewCheckoutHandler(cartRepo, salesRepo, bookRepo)

	reportRepo := repository.NewReportRepo(db)
	reportHandler := handlers.NewReportHandler(reportRepo)

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
	app.Get("/api/customer/orders", middleware.Protect(), customerHandler.GetCustomerOrders)
	app.Get("/api/customer/orders/:id", middleware.Protect(), customerHandler.GetOrderDetails)

	// Cart endpoints
	app.Get("/api/customer/cart", middleware.Protect(), cartHandler.GetCart)
	app.Post("/api/customer/cart/items", middleware.Protect(), cartHandler.AddToCart)
	app.Put("/api/customer/cart/items/:isbn", middleware.Protect(), cartHandler.UpdateCartItem)
	app.Delete("/api/customer/cart/items/:isbn", middleware.Protect(), cartHandler.RemoveFromCart)
	app.Delete("/api/customer/cart", middleware.Protect(), cartHandler.ClearCart)

	// Checkout endpoint
	app.Post("/api/customer/checkout", middleware.Protect(), checkoutHandler.Checkout)

	// Admin-only book management endpoints
	app.Post("/api/admin/books", middleware.AdminOnly(), bookHandler.AddBook)
	app.Put("/api/admin/books/:isbn", middleware.AdminOnly(), bookHandler.UpdateBook)
	app.Patch("/api/admin/books/:isbn/stock", middleware.AdminOnly(), bookHandler.UpdateBookStock)

	// Admin-only publisher order endpoints
	// Note: More specific routes must be registered before less specific ones
	app.Get("/api/admin/publisher-orders/pending", middleware.AdminOnly(), orderHandler.GetPendingPublisherOrders)
	app.Patch("/api/admin/publisher-orders/:id/status", middleware.AdminOnly(), orderHandler.UpdatePublisherOrderStatus)
	app.Put("/api/admin/publisher-orders/:id/confirm", middleware.AdminOnly(), orderHandler.ConfirmPublisherOrder)
	app.Get("/api/admin/publisher-orders/:id", middleware.AdminOnly(), orderHandler.GetPublisherOrder)
	app.Get("/api/admin/publisher-orders", middleware.AdminOnly(), orderHandler.GetAllPublisherOrders)
	app.Post("/api/admin/publisher-orders", middleware.AdminOnly(), orderHandler.PlacePublisherOrder)

	// Admin report endpoints
	app.Get("/api/admin/reports/monthly-sales", middleware.AdminOnly(), reportHandler.GetMonthlySales)
	app.Get("/api/admin/reports/daily-sales", middleware.AdminOnly(), reportHandler.GetDailySales)
	app.Get("/api/admin/reports/top-customers", middleware.AdminOnly(), reportHandler.GetTopCustomers)
	app.Get("/api/admin/reports/top-selling-books", middleware.AdminOnly(), reportHandler.GetTopSellingBooks)
	app.Get("/api/admin/reports/book-orders/:isbn", middleware.AdminOnly(), reportHandler.GetBookOrderCount)

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
