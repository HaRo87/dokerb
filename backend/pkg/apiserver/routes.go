package apiserver

import (
	"github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/haro87/dokerb/docs"
	"github.com/haro87/dokerb/pkg/compute"
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

// GeneralResponse represents a general API response
type GeneralResponse struct {
	Message string `json:"message" example:"ok" format:"string"`
	Route   string `json:"route" example:"/sessions/token" format:"string"`
}

// UsersResponse represents the get users response
type UsersResponse struct {
	Message string   `json:"message" example:"ok" format:"string"`
	Users   []string `json:"users" example:"Tigger,Rabbit" format:"[]string"`
}

// TaskResponse represents the get tasks response
type TaskResponse struct {
	Message string           `json:"message" example:"ok" format:"string"`
	Tasks   []datastore.Task `json:"tasks" format:"[]datastore.Task"`
}

// Task represents a task
type Task struct {
	ID      string `json:"id" example:"TEST01" format:"string"`
	Summary string `json:"summary" example:"a sample task" format:"string"`
}

// Estimate represents a estimate for a task
type Estimate struct {
	Effort            float64 `json:"effort" example:"1.5" format:"float64"`
	StandardDeviation float64 `json:"standarddeviation" example:"0.2" format:"float64"`
}

// CalcEstimate represents the response for calculated average estimate
type CalcEstimate struct {
	Message  string   `json:"message" example:"warning" format:"string"`
	Hint     string   `json:"hint" example:"not all users provided estimates" format:"string"`
	Users    []string `json:"users" example:"Tigger" format:"[]string"`
	Estimate Estimate `json:"estimate" format:"Estimate"`
}

// PerUserEstimate represents a user and task individual estimate
type PerUserEstimate struct {
	TaskID         string  `json:"id" example:"TEST01" format:"string"`
	UserName       string  `json:"user" example:"Tigger" format:"string"`
	BestCase       float64 `json:"b" example:"1.5" format:"float64"`
	MostLikelyCase float64 `json:"m" example:"2.0" format:"float64"`
	WorstCase      float64 `json:"w" example:"3.6" format:"float64"`
}

// PerUserEstimateResponse represents the get estimates response
type PerUserEstimateResponse struct {
	Message   string               `json:"message" example:"ok" format:"string"`
	Estimates []datastore.Estimate `json:"estimates" format:"[]datastore.Estimate"`
}

// User represents a user
type User struct {
	Name string `json:"name" example:"Tigger" format:"string"`
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

	addGetTasksFromSessionRoute(APIGroup, store)

	addAddTaskToSessionRoute(APIGroup, store)

	addRemoveTaskFromSessionRoute(APIGroup, store)

	addUpdateTaskEstimateOfTaskRoute(APIGroup, store)

	addResetEstimateOfTaskRoute(APIGroup, store)

	addAddUserEstimateToSessionRoute(APIGroup, store)

	addRemoveUserEstimateFromSessionRoute(APIGroup, store)

	addGetUserEstimatesFromSessionRoute(APIGroup, store)

	addGetAverageEstimateForTaskFromSessionRoute(APIGroup, store)

	addGetUserWithMaxEstimateDistanceForTaskFromSessionRoute(APIGroup, store)
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
// @Success 200 {object} GeneralResponse
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

		data := GeneralResponse{
			Message: "ok",
			Route:   "/sessions/" + t,
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
// @Param  user body User true "New User"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/users [post]
func addAddUserToSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Post("/sessions/:token/users", func(c *fiber.Ctx) error {
		u := new(User)

		if err := c.BodyParser(u); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.JoinSession(c.Params("token"), u.Name); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
			Route:   "/sessions/" + c.Params("token") + "/users/" + u.Name,
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

// Adding the Get tasks from session route
// @Summary Get the tasks of a session
// @Description Gets all tasks of an existing session
// @Tags task
// @Produce  json
// @Param token path string true "Session Token"
// @Success 200 {object} TasksResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/tasks [get]
func addGetTasksFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Get("/sessions/:token/tasks", func(c *fiber.Ctx) error {
		tasks, e := store.GetTasks(c.Params("token"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := TaskResponse{
			Message: "ok",
			Tasks:   tasks,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Add task to session route
// @Summary Add a new task to a existing session
// @Description Adds a new (non-existing) task to an existing session
// @Tags task
// @Produce  json
// @Param token path string true "Session Token"
// @Param  task body Task true "New Task"
// @Success 200 {object} GeneralResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/tasks [post]
func addAddTaskToSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Post("/sessions/:token/tasks", func(c *fiber.Ctx) error {
		task := new(Task)

		if err := c.BodyParser(task); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.AddTask(c.Params("token"), task.ID, task.Summary); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
			Route:   "/sessions/" + c.Params("token") + "/tasks/" + task.ID,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Remove task from session route
// @Summary Remove a task from a session
// @Description Removes a existing task from an existing session
// @Tags task
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "ID of the task"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/tasks/{id} [delete]
func addRemoveTaskFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token/tasks/:id", func(c *fiber.Ctx) error {
		if err := store.RemoveTask(c.Params("token"), c.Params("id")); err != nil {
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

// Adding the Update estimate of task route
// @Summary Update the estimate of a task
// @Description Updates a estimate of a existing task inside a existing session
// @Tags task
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "ID of the task"
// @Param  estimate body Estimate true "New Estimate"
// @Success 200 {object} GeneralResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/tasks/{id} [put]
func addUpdateTaskEstimateOfTaskRoute(api fiber.Router, store datastore.DataStore) {
	api.Put("/sessions/:token/tasks/:id", func(c *fiber.Ctx) error {
		es := new(Estimate)

		if err := c.BodyParser(es); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		if err := store.AddEstimateToTask(c.Params("token"), c.Params("id"), es.Effort, es.StandardDeviation); err != nil {
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

// Adding the Delete estimate from task route
// @Summary Delete the estimate from a task
// @Description Removes the estimate from an existing task
// @Tags task
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "ID of the task"
// @Success 200 {object} GeneralResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/tasks/{id}/estimate [delete]
func addResetEstimateOfTaskRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token/tasks/:id/estimate", func(c *fiber.Ctx) error {
		if err := store.RemoveEstimateFromTask(c.Params("token"), c.Params("id")); err != nil {
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

// Adding the Add user estimate to session route
// @Summary Add the estimate of a user for a task
// @Description Adds a estimate of a existing user of a existing task inside a existing session
// @Tags estimate
// @Produce  json
// @Param token path string true "Session Token"
// @Param  estimate body PerUserEstimate true "New Estimate"
// @Success 200 {object} GeneralResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/estimates [post]
func addAddUserEstimateToSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Post("/sessions/:token/estimates", func(c *fiber.Ctx) error {
		es := new(PerUserEstimate)

		if err := c.BodyParser(es); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(400).JSON(data)
		}

		est := datastore.Estimate{
			TaskID:         es.TaskID,
			UserName:       es.UserName,
			BestCase:       es.BestCase,
			MostLikelyCase: es.MostLikelyCase,
			WorstCase:      es.WorstCase,
		}

		if err := store.AddEstimate(c.Params("token"), est); err != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  err.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := GeneralResponse{
			Message: "ok",
			Route:   "/sessions/" + c.Params("token") + "/estimates/" + es.UserName + "/" + es.TaskID,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Remove user estimate from session route
// @Summary Remove the estimate of a user for a task
// @Description Removes a estimate of a existing user of a existing task inside a existing session
// @Tags estimate
// @Produce  json
// @Param token path string true "Session Token"
// @Param  user path string true "User Name"
// @Param  id path string true "Task ID"
// @Success 200 {object} GeneralResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/estimates/{user}/{id} [delete]
func addRemoveUserEstimateFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Delete("/sessions/:token/estimates/:user/:id", func(c *fiber.Ctx) error {

		est := datastore.Estimate{
			TaskID:   c.Params("id"),
			UserName: c.Params("user"),
		}

		if err := store.RemoveEstimate(c.Params("token"), est); err != nil {
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

// Adding the Get user estimates from session route
// @Summary Get the estimates of all users for all tasks
// @Description Gets all estimates of all existing users of all existing tasks inside a existing session
// @Tags estimate
// @Produce  json
// @Param token path string true "Session Token"
// @Success 200 {object} PerUserEstimateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/estimates [get]
func addGetUserEstimatesFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Get("/sessions/:token/estimates", func(c *fiber.Ctx) error {

		ests, e := store.GetEstimates(c.Params("token"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := PerUserEstimateResponse{
			Message:   "ok",
			Estimates: ests,
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Get average user estimate from session route
// @Summary Get the average estimate of all users for a specific task
// @Description Gets the average estimate of all existing users of a existing task inside a existing session
// @Tags estimate
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "Task ID"
// @Success 200 {object} CalcEstimate
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/estimates/{id} [get]
func addGetAverageEstimateForTaskFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Get("/sessions/:token/estimates/:id", func(c *fiber.Ctx) error {

		ests, e := store.GetEstimates(c.Params("token"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		ests, e = compute.ExtractEstimatesForTask(ests, c.Params("id"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		users, ue := store.GetUsers(c.Params("token"))

		if ue != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  ue.Error(),
			}
			return c.Status(500).JSON(data)
		}

		avge, ae := compute.CalculateAverageEstimate(ests, c.Params("id"))

		if ae != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  ae.Error(),
			}
			return c.Status(500).JSON(data)
		}

		message := "ok"
		hint := ""

		for _, es := range ests {
			users = checkForAllUsers(users, es.UserName)
		}

		if len(users) > 0 {
			message = "warning"
			hint = "not all users did provide estimates"
		}

		data := CalcEstimate{
			Message: message,
			Hint:    hint,
			Users:   users,
			Estimate: Estimate{
				Effort:            avge.GetEffort(),
				StandardDeviation: avge.GetStandardDeviation(),
			},
		}
		return c.Status(200).JSON(data)
	})
}

// Adding the Get max distance users for estimate from session route
// @Summary Get the users with max distance between their estimates for a specific task
// @Description Gets the users with max distance in their estimates of a existing task inside a existing session
// @Tags estimate
// @Produce  json
// @Param token path string true "Session Token"
// @Param id path string true "Task ID"
// @Success 200 {object} UsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{token}/estimates/{id}/users/distance [get]
func addGetUserWithMaxEstimateDistanceForTaskFromSessionRoute(api fiber.Router, store datastore.DataStore) {
	api.Get("/sessions/:token/estimates/:id/users/distance", func(c *fiber.Ctx) error {

		ests, e := store.GetEstimates(c.Params("token"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		ests, e = compute.ExtractEstimatesForTask(ests, c.Params("id"))

		if e != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  e.Error(),
			}
			return c.Status(500).JSON(data)
		}

		users, ae := compute.GetUsersWithMaxDistanceBetweenEffort(ests, c.Params("id"))

		if ae != nil {
			data := ErrorResponse{
				Message: "error",
				Reason:  ae.Error(),
			}
			return c.Status(500).JSON(data)
		}

		data := UsersResponse{
			Message: "ok",
			Users:   users,
		}
		return c.Status(200).JSON(data)
	})
}

func checkForAllUsers(users []string, user string) []string {
	for i, u := range users {
		if u == user {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	return users
}
