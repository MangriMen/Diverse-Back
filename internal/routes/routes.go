package routes

import (
	"path/filepath"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/users", controllers.GetUsers)
	route.Get("/users/:id", controllers.GetUser)

	route.Static("data/image/", filepath.Join(configs.DataPath, configs.MIMEBaseImage), fiber.Static{Compress: true})

	route.Post("/login", controllers.LoginUser)
	route.Post("/register", controllers.CreateUser)
}

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/fetch", middleware.JWTProtected(), controllers.FetchUser)

	route.Patch("/users/:id", middleware.JWTProtected(), controllers.UpdateUser)
	route.Delete("/users/:id", middleware.JWTProtected(), controllers.DeleteUser)

	route.Post("/data", middleware.JWTProtected(), controllers.UploadData)
}
