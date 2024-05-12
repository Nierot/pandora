package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/nierot/pandora/models"
	"github.com/nierot/pandora/routes"
	"github.com/spf13/viper"
)

func main() {
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		CaseSensitive: false,
		Views: engine,
	})

	app.Static("/static", "./public")

	config()

	models.InitDB()
	
	routes.InitPublic(app)
	routes.InitAuth(app)
	routes.InitProtected(app)
	
	app.Listen(":3000")
}

func config() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("AUTH_USERNAME", "admin")
	viper.SetDefault("AUTH_PASSWORD", "123456")
	viper.SetDefault("DATABASE_URL", "root:@tcp(localhost:3306)/pandora")

	viper.SafeWriteConfig()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

