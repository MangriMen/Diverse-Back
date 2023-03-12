// Package routes provides routes aggregation for application
package routes

import (
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes sets up public routes for an API version X
// by defining a group of routes under the prefix "/api/vX".
func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	UserPublicRoutes(route)
}

// PrivateRoutes sets up private routes for an API version X
// by defining a group of routes under the prefix "/api/vX".
func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	UserPrivateRoutes(route)
	PostPrivateRoutes(route)
	DataPrivateRoutes(route)
}
