package apiserver

import (
	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
	logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/haro87/dokerb/pkg/datastore"
	"os"
)

// APIServer struct
type APIServer struct {
	config *Config
	ds     datastore.DataStore
}

// NewServer method for init new server instance
func NewServer(config *Config, ds datastore.DataStore) *APIServer {
	return &APIServer{
		config: config,
		ds:     ds,
	}
}

// Start method for start new server
func (s *APIServer) Start() *fiber.App {
	// Initialize a new app
	app := fiber.New()

	// Register middlewares
	app.Use(
		cors.New(), // Add CORS to each route
		// Simple logger
		logger.New(
			logger.Config{
				Format:     "${time} [${status}] ${method} ${path} (${latency})\n",
				TimeFormat: "Mon, 2 Jan 2006 15:04:05 MST",
				Output:     os.Stdout,
			},
		),
	)

	// Add static files, if prefix and path was defined in config
	if s.config.Static.Prefix != "" && s.config.Static.Path != "" {
		app.Static(s.config.Static.Prefix, s.config.Static.Path)
	}

	// Register API routes
	Routes(app, s.ds, s.config)

	return app
}
