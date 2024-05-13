package main

import (
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/nierot/pandora/models"
	"github.com/nierot/pandora/routes"
	"github.com/spf13/viper"

	"github.com/gomarkdown/markdown"
	mdHtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	engine := html.New("./templates", ".html")
	engine.AddFunc("markdown", func(s string) template.HTML {

		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse([]byte(s))
		
		htmlFlags := mdHtml.CommonFlags | mdHtml.HrefTargetBlank
		opts := mdHtml.RendererOptions{Flags: htmlFlags}
		renderer := mdHtml.NewRenderer(opts)
		
		
		md := markdown.Render(doc, renderer)

		return template.HTML(md)
	})

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
	viper.AddConfigPath("/etc/pandora/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	// viper.SetDefault("AUTH_USERNAME", "admin")
	// viper.SetDefault("AUTH_PASSWORD", "123456")
	// viper.SetDefault("DATABASE_URL", "root:@tcp(localhost:3306)/pandora")
	// viper.SafeWriteConfig()

	test := viper.GetString("DATABASE_URL")

	fmt.Println("database test: ", test)

	if test == "" {
		panic("DATABASE_URL is not set")
	}

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

