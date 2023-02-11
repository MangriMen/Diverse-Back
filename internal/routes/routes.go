package routes

import (
	"github.com/MangriMen/Value-Back/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/users", controllers.GetUsers)
	route.Get("/users/:id", controllers.GetUser)
	// route.Get("/token/new", controllers.GetNewAccessToken)
}

func PrivateRoutes(a *fiber.App) {
	// route := a.Group("/api/v1")

	// route.Post("/user", middleware.JWTProtected(), controllers.CreateUser)

	// route.Patch("/user", middleware.JWTProtected(), controllers.UpadteUser)

	// route.Delete("/user", middleware.JWTProtected(), controllers.DeleteUser)
}
