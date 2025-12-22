package main

import (
	"log"

	"bookstore-project/internal/config"
	"bookstore-project/internal/database"
	"bookstore-project/internal/handlers"
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

	app := fiber.New()

	Repo := repository.NewBookRepo(db)
	handler := handlers.NewBookHandler(Repo)

	app.Get("/books", handler.GetAllBooks)
	app.Get("/books/:isbn",handler.GetBookByIsbn)
	app.Get("/books/author/:name",handler.GetBookByAuthor)

	app.Listen(":3001")

}
