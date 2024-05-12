package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/spf13/viper"
)

func InitAuth(app *fiber.App) {
	username := viper.GetString("AUTH_USERNAME")
	password := viper.GetString("AUTH_PASSWORD")

	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
				username: password,
		},
	}))
}