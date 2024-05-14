package routes

import (
	"fmt"
	"strconv"

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

	app.Post("/nieuw", func (c *fiber.Ctx) error {
		player := c.FormValue("player")
		reason := c.FormValue("reason")

		pid, err := strconv.Atoi(player)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		err = models.AddBak(pid, reason)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.RedirectBack("/nieuw")
	})

	app.Get("/nieuwe-blog", func (c *fiber.Ctx) error {
		// if ?edit=true then show edit form

		edit := c.Query("edit")

		if edit == "true" {
			id, err := strconv.Atoi(c.Query("id"))

			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			blog, err := models.GetBlogEntry(id)

			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			players, err := models.GetPlayers()

			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			return c.Render("edit_blog", fiber.Map{
				"Players": players,
				"Blog": blog,
			}, "layouts/base")
		}


		players, err := models.GetPlayers()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Render("add_blog", fiber.Map{
			"Players": players,	
		}, "layouts/base")
	})

	app.Post("/nieuwe-blog", func (c *fiber.Ctx) error {
		player := c.FormValue("player")
		title := c.FormValue("title")
		content := c.FormValue("content")

		pid, err := strconv.Atoi(player)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		idx, err := models.AddBlog(pid, title, content)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Redirect("/blog/" + strconv.Itoa(idx))
	})

	app.Post("/edit-blog", func (c *fiber.Ctx) error {
		id := c.FormValue("id")
		player := c.FormValue("player")
		title := c.FormValue("title")
		content := c.FormValue("content")

		fmt.Println("id " + id)

		pid, err := strconv.Atoi(player)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		nid, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println(err)
			return c.Status(500).SendString(err.Error())
		}

		err = models.EditBlog(nid, pid, title, content)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Redirect("/blog/" + id)
	})
}