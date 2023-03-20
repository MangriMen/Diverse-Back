package routes

import (
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	UserPublicRoutes(route)
}

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	UserPrivateRoutes(route)
	PostPrivateRoutes(route)
}
