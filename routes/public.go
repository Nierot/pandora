package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nierot/pandora/models"
)

func InitPublic(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		bakken, err := models.GetAmountOfBakkenPerPlayer()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("index", fiber.Map{
			"Bakken": bakken,
		}, "layouts/base")
	})

	app.Get("/bakken", func (c *fiber.Ctx) error {
		bakken, err := models.GetBakken()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("bakken", fiber.Map{
			"Bakken": bakken,
		}, "layouts/base")
	})
}