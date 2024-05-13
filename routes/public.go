package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nierot/pandora/models"
)

func InitPublic(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		bakken, err := models.GetAmountOfBakkenPerPlayer()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		blog, err := models.GetLatestBlogEntry()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		blogList, err := models.GetLast3BlogEntries()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("index", fiber.Map{
			"Bakken": bakken,
			"Blog": &blog,
			"BlogList": blogList,
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

	app.Get("/blog", func (c *fiber.Ctx) error {
		blogs, err := models.GetBlogEntries()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("blog_index", fiber.Map{
			"Blogs": blogs,
		}, "layouts/base")
	})


	app.Get("/blog/:id", func (c *fiber.Ctx) error {
		id := c.Params("id")

		nid, err := strconv.Atoi(id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		blog, err := models.GetBlogEntry(nid)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		previous := blog.ID - 1

		if previous < 1 {
			previous = 1
		}

		nextBlog, _ := models.GetBlogEntry(blog.ID + 1)

		next := 0
		hasNext := false

		if nextBlog == nil {
			next = 0
		} else {
			hasNext = true
			next = nextBlog.ID
		}

		return c.Render("blog", fiber.Map{
			"Blog": blog,
			"Previous": previous,
			"Next": next,
			"HasNext": hasNext,
		}, "layouts/base")
	})
}