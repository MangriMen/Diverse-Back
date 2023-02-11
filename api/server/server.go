package server

import (
	"github.com/MangriMen/Value-Back/internal/helpers"
	"github.com/MangriMen/Value-Back/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitApi() {
	app := fiber.New()

	app.Use(
		cors.New(),
		logger.New(),
	)

	routes.PublicRoutes(app)
	// routes.PrivateRoutes(app)

	helpers.StartServerWithGracefulShutdown(app)
}
