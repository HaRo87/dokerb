package apiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haro87/dokerb/pkg/datastore"
)

// WorkPackage represents a work package
type WorkPackage struct {
	ID      string
	Summary string
}

// Estimate represents a work package estimate
type Estimate struct {
	Effort            float64
	StandardDeviation float64
}

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

	APIGroup.Delete("/sessions/:token", func(c *fiber.Ctx) error {
		if err := store.RemoveSession(c.Params("token")); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Post("/sessions/:token/users/:name", func(c *fiber.Ctx) error {
		if err := store.JoinSession(c.Params("token"), c.Params("name")); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
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

	APIGroup.Delete("/sessions/:token/users/:name", func(c *fiber.Ctx) error {
		if err := store.LeaveSession(c.Params("token"), c.Params("name")); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Get("/sessions/:token/workpackages", func(c *fiber.Ctx) error {
		wps, e := store.GetWorkPackages(c.Params("token"))

		var data fiber.Map

		if e != nil {
			data = fiber.Map{
				"message": "error",
				"reason":  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data = fiber.Map{
			"message":      "ok",
			"workpackages": wps,
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Post("/sessions/:token/workpackages", func(c *fiber.Ctx) error {
		wp := new(WorkPackage)

		if err := c.BodyParser(wp); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.AddWorkPackage(c.Params("token"), wp.ID, wp.Summary); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Delete("/sessions/:token/workpackages/:id", func(c *fiber.Ctx) error {
		if err := store.RemoveWorkPackage(c.Params("token"), c.Params("id")); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Put("/sessions/:token/workpackages/:id", func(c *fiber.Ctx) error {
		es := new(Estimate)

		if err := c.BodyParser(es); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.AddEstimate(c.Params("token"), c.Params("id"), es.Effort, es.StandardDeviation); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})

	APIGroup.Delete("/sessions/:token/workpackages/:id/estimate", func(c *fiber.Ctx) error {
		if err := store.RemoveEstimate(c.Params("token"), c.Params("id")); err != nil {
			data := fiber.Map{
				"message": "error",
				"reason":  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := fiber.Map{
			"message": "ok",
		}
		return c.Status(200).JSON(data)
	})

}
