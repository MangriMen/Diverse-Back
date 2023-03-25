// Package middleware provides middlewares for the application
package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// FiberMiddleware sets up the CORS configration and logger.
func FiberMiddleware(a *fiber.App) {
	a.Use(
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowOrigins:     os.Getenv("FRONTEND_BASE_URL"),
			AllowHeaders: `Origin,
				Content-Type,
				Accept,
				Content-Length,
				Accept-Language,
				Accept-Encoding,
				Connection,
				Authorization`,
			AllowMethods: "POST, GET, PATCH, DELETE, PUT"}),

		logger.New(),
	)
}
