package routes

import (
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func PostPrivateRoutes(route fiber.Router) {
	route.Get("/posts", middleware.JWTProtected(), controllers.GetPosts)

	route.Get("/posts/:id", middleware.JWTProtected(), controllers.GetPost)

	route.Post("/posts", middleware.JWTProtected(), controllers.CreatePost)

	route.Patch("/posts/:id", middleware.JWTProtected(), controllers.UpdatePost)

	route.Delete("/posts/:id", middleware.JWTProtected(), controllers.DeletePost)

	PostCommentPrivateRoutes(route)
}

func PostCommentPrivateRoutes(route fiber.Router) {
	posts := route.Group("/posts/:id")

	posts.Post("/comments", middleware.JWTProtected(), controllers.AddComment)

	posts.Patch("/comments/:comment", middleware.JWTProtected(), controllers.UpdateComment)

	posts.Delete("/comments/:comment", middleware.JWTProtected(), controllers.DeleteComment)
}
