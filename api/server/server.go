// Package server provides functionality to run REST API server
package server

import (
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/MangriMen/Diverse-Back/internal/routes"
	"github.com/gofiber/fiber/v2"
)

// InitAPI is used for initialize a new instance of fiber web application.
func InitAPI() {
	config := configs.FiberConfig()
	app := fiber.New(config)

	middleware.FiberMiddleware(app)
	middleware.Compress(app)

	routes.PublicRoutes(app)
	routes.PrivateRoutes(app)

	helpers.StartServerWithGracefulShutdown(app)
}
