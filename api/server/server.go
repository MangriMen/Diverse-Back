package server

import (
	"github.com/MangriMen/Value-Back/configs"
	"github.com/MangriMen/Value-Back/internal/helpers"
	"github.com/MangriMen/Value-Back/internal/middleware"
	"github.com/MangriMen/Value-Back/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func InitApi() {
	config := configs.FiberConfig()
	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	routes.PublicRoutes(app)
	routes.PrivateRoutes(app)

	helpers.StartServerWithGracefulShutdown(app)
}
