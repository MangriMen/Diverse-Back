package routes

import (
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserPublicRoutes(route fiber.Router) {
	route.Get("/users", controllers.GetUsers)
	route.Get("/users/:id", controllers.GetUser)

	route.Post("/login", controllers.LoginUser)
	route.Post("/register", controllers.CreateUser)
}

func UserPrivateRoutes(route fiber.Router) {
	route.Get("/fetch", middleware.JWTProtected(), controllers.FetchUser)

	route.Patch("/users/:id", middleware.JWTProtected(), controllers.UpdateUser)

	route.Delete("/users/:id", middleware.JWTProtected(), controllers.DeleteUser)
}
