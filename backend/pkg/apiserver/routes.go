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
					"url":  "https://create-go.app/",
				},
				{
					"name": "Detailed guides",
					"url":  "https://create-go.app/detailed-guides/",
				},
				{
					"name": "GitHub",
					"url":  "https://github.com/create-go-app/cli",
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
				"message": e,
			}
			return c.Status(500).JSON(data)
		}

		data = fiber.Map{
			"message": "ok",
			"token":   t,
		}
		return c.Status(200).JSON(data)
	})
}
