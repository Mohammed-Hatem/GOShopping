package main

import (
	"log"

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

	Repo := repository.NewBookRepo(db)
	handler := handlers.NewBookHandler(Repo)

	customerRepo := repository.NewCustomerRepo(db)
	customerHandler := handlers.NewCustomerHandler(customerRepo)

	cartRepo := repository.NewCartRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	cartHandler := handlers.NewCartHandler(cartRepo, orderRepo)

	app.Get("/books", handler.GetAllBooks)
	app.Get("/books/author/:name", handler.GetBookByAuthor)
	app.Get("/books/pubyear/:year", handler.GetBookByPubyear)
	app.Get("/books/title/:title", handler.GetBookByTitle)
	app.Get("/books/category/:category", handler.GetBookByCategory)
	app.Get("/books/:isbn", handler.GetBookByIsbn)

	app.Post("/signup", customerHandler.Signup)
	app.Post("/login", customerHandler.Login)
	app.Get("/profile", middleware.Protect(), customerHandler.GetProfile)
	app.Patch("/profile", middleware.Protect(), customerHandler.UpdateProfile)

	// Cart + checkout (Customer)
	app.Post("/cart/items", middleware.Protect(), cartHandler.AddToCart)
	app.Get("/cart", middleware.Protect(), cartHandler.ViewCart)
	app.Delete("/cart/items/:isbn", middleware.Protect(), cartHandler.RemoveItem)
	app.Post("/cart/checkout", middleware.Protect(), cartHandler.Checkout)

	err = app.Listen(":3001")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}



}
