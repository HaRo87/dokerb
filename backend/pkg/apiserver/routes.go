package apiserver

import (
	"github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/haro87/dokerb/docs"
	"github.com/haro87/dokerb/pkg/datastore"
)

// DocEntry represents a single documentation entry
type DocEntry struct {
	Name string `json:"name" example:"GitHub" format:"string"`
	URL  string `json:"url" example:"https://github.com/HaRo87" format:"string"`
}

// DocResponse represents the full documentation response
type DocResponse struct {
	Message string     `json:"message" example:"ok" format:"string"`
	Results []DocEntry `json:"results"`
}

// ErrorResponse represents a response in case an error ocurred
type ErrorResponse struct {
	Message string `json:"message" example:"error" format:"string"`
	Reason  string `json:"reason" example:"oops, something went wrong" format:"string"`
}

// SessionResponse represents the response of the create session call
type SessionResponse struct {
	Message string `json:"message" example:"ok" format:"string"`
	Token   string `json:"token" example:"12345678901234567890123456789012" format:"string"`
}

// GeneralResponse represents a general API response
type GeneralResponse struct {
	Message string `json:"message" example:"ok" format:"string"`
}

// UsersResponse represents the get users response
type UsersResponse struct {
	Message string   `json:"message" example:"ok" format:"string"`
	Users   []string `json:"users" example:"Tigger,Rabbit" format:"[]string"`
}

// WorkPackagesResponse represents the get work packages response
type WorkPackagesResponse struct {
	Message      string                  `json:"message" example:"ok" format:"string"`
	Workpackages []datastore.WorkPackage `json:"workpackages" format:"[]datastore.WorkPackage"`
}

// WorkPackage represents a work package
type WorkPackage struct {
	ID      string `json:"id" example:"TEST01" format:"string"`
	Summary string `json:"summary" example:"a sample task" format:"string"`
}

// Estimate represents a work package estimate
type Estimate struct {
	Effort            float64 `json:"effort" example:"1.5" format:"float64"`
	StandardDeviation float64 `json:"standarddeviation" example:"0.2" format:"float64"`
}

// Routes list of the available routes for project
// @title Doker Backend API
// @version 0.1.0
// @description A backend for playing Planning Poker with Delphi estimate method.

// @contact.name HaRo87
// @contact.url https://github.com/HaRo87

// @license.name MIT
// @license.url https://github.com/HaRo87/dokerb/blob/main/LICENSE

// @host localhost:5000
// @BasePath /api
func Routes(app *fiber.App, store datastore.DataStore) {
	// Create group for API routes
	APIGroup := app.Group("/api")

	APIGroup.Get("/swagger/*", swagger.Handler)

	addDocRoute(APIGroup)

	addCreateSessionRoute(APIGroup, store)

	addRemoveSessionRoute(APIGroup, store)

	addAddUserToSessionRoute(APIGroup, store)

	addGetUsersFromSessionRoute(APIGroup, store)

	addRemoveUserFromSessionRoute(APIGroup, store)

	addGetWorkPackagesFromSessionRoute(APIGroup, store)

	addAddWorkPackageToSessionRoute(APIGroup, store)

	addRemoveWorkPackageFromSessionRoute(APIGroup, store)

	addUpdateWorkPackageEstimateOfWorkPackageRoute(APIGroup, store)

	addResetEstimateOfWorkPackageRoute(APIGroup, store)

}

// Adding the documentation route
// @Summary Get the documentation info
// @Description Get a list of helpful documentation resources
// @Tags documentation
// @Produce  json
// @Success 200 {object} DocResponse
// @Router /docs [get]
func addDocRoute(api fiber.Router) {

	api.Get("/docs", func(c *fiber.Ctx) error {
		// Set JSON data
		data := DocResponse{
			Message: "ok",
			Results: []DocEntry{
				{
					Name: "Documentation",
					URL:  "https://haro87.github.io/doker-meta",
				},
				{
					Name: "GitHub",
					URL:  "https://github.com/HaRo87/dokerb",
				},
				{
					Name: "Swagger",
					URL:  "/api/swagger/",
				},
			},
		}

		// Set 200 OK status and return JSON
		return c.Status(200).JSON(data)
	})
}

// Adding the create session route
// @Summary Create a new Doker session
// @Description Creates a new Doker session and responds with the corresponding token
// @Tags session
// @Produce  json
// @Success 200 {object} SessionResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions [post]
func addCreateSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Post("/sessions", func(c *fiber.Ctx) error {
		t, err := store.CreateSession()

		if err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := SessionResponse{
			Message: "ok",
			Token:   t,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the remove session route
// @Summary Delete a existing Doker session
// @Description Deletes a existing Doker session based on the provided token
// @Tags session
// @Produce  json
// @Param token path string true "Session Token"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token} [delete]
func addRemoveSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token", func(c *fiber.Ctx) error {
		if err := store.RemoveSession(c.Params("token")); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Add user to session route
// @Summary Add a new user to a existing session
// @Description Adds a new (non-existing) user to an existing session
// @Tags user
// @Produce  json
// @Param token path string true "Session Token"
// @Param name path string true "Name of the user"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/users/{name} [post]
func addAddUserToSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Post("/sessions/:token/users/:name", func(c *fiber.Ctx) error {
		if err := store.JoinSession(c.Params("token"), c.Params("name")); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Get users from session route
// @Summary Get the users of an existing session
// @Description Gets all users of an existing session
// @Tags user
// @Produce  json
// @Param token path string true "Session Token"
// @Success 200 {object} UsersResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/users [get]
func addGetUsersFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Get("/sessions/:token/users", func(c *fiber.Ctx) error {
		u, e := store.GetUsers(c.Params("token"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := UsersResponse{
			Message: "ok",
			Users:   u,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Remove user from session route
// @Summary Remove a user from a session
// @Description Removes a existing user from an existing session
// @Tags user
// @Produce  json
// @Param token path string true "Session Token"
// @Param name path string true "Name of the user"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/users/{name} [delete]
func addRemoveUserFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token/users/:name", func(c *fiber.Ctx) error {
		if err := store.LeaveSession(c.Params("token"), c.Params("name")); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Get work packages from session route
// @Summary Get the work packages of a session
// @Description Gets all work packages of an existing session
// @Tags workpackage
// @Produce  json
// @Param token path string true "Session Token"
// @Success 200 {object} WorkPackagesResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/workpackages [get]
func addGetWorkPackagesFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Get("/sessions/:token/workpackages", func(c *fiber.Ctx) error {
		wps, e := store.GetWorkPackages(c.Params("token"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := WorkPackagesResponse{
			Message:      "ok",
			Workpackages: wps,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Add work package to session route
// @Summary Add a new work package to a existing session
// @Description Adds a new (non-existing) work package to an existing session
// @Tags workpackage
// @Produce  json
// @Param token path string true "Session Token"
// @Param  workpackage body WorkPackage true "New Work Package"
// @Success 200 {object} GeneralResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/workpackages [post]
func addAddWorkPackageToSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Post("/sessions/:token/workpackages", func(c *fiber.Ctx) error {
		wp := new(WorkPackage)

		if err := c.BodyParser(wp); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.AddWorkPackage(c.Params("token"), wp.ID, wp.Summary); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Remove work package from session route
// @Summary Remove a work package from a session
// @Description Removes a existing work package from an existing session
// @Tags workpackage
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "ID of the work package"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/workpackages/{id} [delete]
func addRemoveWorkPackageFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token/workpackages/:id", func(c *fiber.Ctx) error {
		if err := store.RemoveWorkPackage(c.Params("token"), c.Params("id")); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Update estimate of work package route
// @Summary Update the estimate of a work package
// @Description Updates a estimate of a existing work package inside a existing session
// @Tags workpackage
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "ID of the work package"
// @Param  estimate body Estimate true "New Estimate"
// @Success 200 {object} GeneralResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/workpackages/{id} [put]
func addUpdateWorkPackageEstimateOfWorkPackageRoute(api fiber.Router, store datastore.DataStore) {
	api.Put("/sessions/:token/workpackages/:id", func(c *fiber.Ctx) error {
		es := new(Estimate)

		if err := c.BodyParser(es); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.AddEstimate(c.Params("token"), c.Params("id"), es.Effort, es.StandardDeviation); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Delete estimate from work package route
// @Summary Delete the estimate from a work package
// @Description Removes the estimate from an existing work package
// @Tags workpackage
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "ID of the work package"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/workpackages/{id}/estimate [delete]
func addResetEstimateOfWorkPackageRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token/workpackages/:id/estimate", func(c *fiber.Ctx) error {
		if err := store.RemoveEstimate(c.Params("token"), c.Params("id")); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
		}
		return c.Status(200).JSON(data)
	})
}
