package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nierot/pandora/models"
)

func InitProtected(app *fiber.App) {
	app.Get("/nieuw", func (c *fiber.Ctx) error {
		players, err := models.GetPlayers()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("add_bak", fiber.Map{
			"Players": players,	
		}, "layouts/base")
	})
}