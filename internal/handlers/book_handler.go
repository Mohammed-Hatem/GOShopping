package handlers


import (
	"net/http"
	"github.com/gofiber/fiber/v2"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {

	app := fiber.New()

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")

	defer app.Shutdown()
}
