package apiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haro87/dokerb/pkg/datastore"
)

// Routes list of the available routes for project
func Routes(app *fiber.App, store datastore.DataStore) {
	// Create group for API routes
	APIGroup := app.Group("/api")

	// API routes
	APIGroup.Get("/docs", func(c *fiber.Ctx) error {
		// Set JSON data
		data := fiber.Map{
			"message": "ok",
			"results": []fiber.Map{
				{
					"name": "Documentation",
					"url":  "https://haro87.github.io/doker-meta",
				},
				{
					"name": "GitHub",
					"url":  "https://github.com/HaRo87/dokerb",
				},
			},
		}

		// Set 200 OK status and return JSON
		return c.Status(200).JSON(data)
	})

	APIGroup.Post("/sessions", func(c *fiber.Ctx) error {
		t, e := store.CreateSession()

		var data fiber.Map

		if e != nil {
			data = fiber.Map{
				"message": "error",
				"reason":  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data = fiber.Map{
			"message": "ok",
			"token":   t,
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Get("/sessions/:token/users", func(c *fiber.Ctx) error {
		u, e := store.GetUsers(c.Params("token"))

		var data fiber.Map

		if e != nil {
			data = fiber.Map{
				"message": "error",
				"reason":  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data = fiber.Map{
			"message": "ok",
			"users":   u,
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Post("/sessions/:token/users/:name", func(c *fiber.Ctx) error {
		e := store.JoinSession(c.Params("token"), c.Params("name"))

		var data fiber.Map

		if e != nil {
			data = fiber.Map{
				"message": "error",
				"reason":  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data = fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})
}
