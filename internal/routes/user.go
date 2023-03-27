package routes

import (
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// UserPublicRoutes sets up the public routes for user-related API endpoints
// such as getting a list of users, getting a specific user by ID,
// logging in a user, and register a new user.
func UserPublicRoutes(route fiber.Router) {
	route.Get("/users", controllers.GetUsers)
	route.Get("/users/:id", controllers.GetUser)

	route.Post("/login", controllers.LoginUser)
	route.Post("/register", controllers.CreateUser)
}

// UserPrivateRoutes sets up private routes for authenticated users.
// These routes require a valid JWT for authentication and authorization to access the endpoints.
// It includes endpoints for fetching and updating user information, as well as deleting user accounts.
func UserPrivateRoutes(route fiber.Router) {
	route.Get("/fetch", middleware.JWTProtected(), controllers.FetchUser)

	route.Patch("/users/:id", middleware.JWTProtected(), controllers.UpdateUser)

	route.Delete("/users/:id", middleware.JWTProtected(), controllers.DeleteUser)
}
