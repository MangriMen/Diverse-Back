package routes

import (
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Post("/login", controllers.LoginUser)
	route.Post("/register", controllers.CreateUser)

	route.Get("/users", controllers.GetUsers)
	route.Get("/users/:id", controllers.GetUser)
}

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Patch("/user", middleware.JWTProtected(), controllers.UpdateUser)
	route.Delete("/user", middleware.JWTProtected(), controllers.DeleteUser)
}
