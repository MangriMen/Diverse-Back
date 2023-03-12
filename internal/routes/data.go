package routes

import (
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// DataPrivateRoutes sets up private routes for authenticated users.
// These routes require a valid JWT for authentication and authorization to access the endpoints.
// It includes endpoints for uploading data like images.
func DataPrivateRoutes(route fiber.Router) {
	route.Post("/data", middleware.JWTProtected(), controllers.UploadData)
}
