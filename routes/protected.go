package routes

import (
	"github.com/gofiber/fiber/v2"
)

func InitProtected(app *fiber.App) {
	app.Get("/nieuw", func(c *fiber.Ctx) error {
		return c.SendString("This is a protected route")
	})
}