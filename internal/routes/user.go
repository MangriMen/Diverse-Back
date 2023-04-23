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
	route.Get("/users/:user", controllers.GetUser)

	route.Post("/login", controllers.LoginUser)
	route.Post("/register", controllers.CreateUser)
}

// UserPrivateRoutes sets up private routes for authenticated users.
// These routes require a valid JWT for authentication and authorization to access the endpoints.
// It includes endpoints for fetching and updating user information, as well as deleting user accounts.
func UserPrivateRoutes(route fiber.Router) {
	route.Get("/fetch", middleware.JWTProtected(), controllers.FetchUser)

	route.Patch("/users/:user", middleware.JWTProtected(), controllers.UpdateUser)

	route.Delete("/users/:user", middleware.JWTProtected(), controllers.DeleteUser)

	UserRelationPrivateRoutes(route)
}

// UserRelationPrivateRoutes sets up private routes for authenticated users.
// These routes require a valid JWT for authentication and authorization to access the endpoints.
// It includes endpoints for fetching, updating, removing user relations.
func UserRelationPrivateRoutes(route fiber.Router) {
	users := route.Group("/users/:user")

	users.Get("/relations/count", middleware.JWTProtected(), controllers.GetRelationsCount)

	users.Get("/relations", middleware.JWTProtected(), controllers.GetRelations)

	users.Get("/relations/:relationUser", middleware.JWTProtected(), controllers.GetRelationStatus)

	users.Post("/relations/:relationUser", middleware.JWTProtected(), controllers.AddRelation)

	users.Delete("/relations/:relationUser", middleware.JWTProtected(), controllers.DeleteRelation)
}
