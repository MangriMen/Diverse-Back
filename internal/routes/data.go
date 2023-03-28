package routes

import (
	"path/filepath"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// DataPublicRoutes sets up the public routes for data-related API endpoints
// such as getting a image.
func DataPublicRoutes(route fiber.Router) {
	// swagger:route GET /data/image/raw/{image} Data getImageRaw
	// Returns the image by id
	//
	// Produces:
	//   - image/webp
	//   - application/json
	//
	// Schemes: http, https
	//
	// Responses:
	//   200: StatusOk
	//   default: ErrorResponse

	route.Static(
		"data/image/raw",
		filepath.Join(configs.DataPath, configs.MIMEBaseImage),
		fiber.Static{Compress: true},
	)
}

// DataPrivateRoutes sets up private routes for authenticated users.
// These routes require a valid JWT for authentication and authorization to access the endpoints.
// It includes endpoints for uploading data like images.
func DataPrivateRoutes(route fiber.Router) {
	route.Post("/data", middleware.JWTProtected(), controllers.UploadData)
}
