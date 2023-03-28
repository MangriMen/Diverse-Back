package routes

import (
	"github.com/MangriMen/Diverse-Back/internal/controllers"
	"github.com/MangriMen/Diverse-Back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// PostPrivateRoutes sets up the private routes for the post-related endpoints,
// which require JWT authentication to access. It includes routes
// for retrieving, creating, updating, and deleting posts.
// Additionally, it sets up the private routes for post comments.
func PostPrivateRoutes(route fiber.Router) {
	route.Get("/posts", middleware.JWTProtected(), controllers.GetPosts)

	route.Get("/posts/:post", middleware.JWTProtected(), controllers.GetPost)

	route.Post("/posts", middleware.JWTProtected(), controllers.CreatePost)

	route.Post("/posts/:post/like", middleware.JWTProtected(), controllers.LikePost)

	route.Patch("/posts/:post", middleware.JWTProtected(), controllers.UpdatePost)

	route.Delete("/posts/:post", middleware.JWTProtected(), controllers.DeletePost)

	route.Delete("/posts/:post/like", middleware.JWTProtected(), controllers.UnlikePost)

	PostCommentPrivateRoutes(route)
}

// PostCommentPrivateRoutes sets up the private routes for the comment-related endpoints,
// which require JWT authentication to access. It includes routes
// for retrieving, creating, updating, and deleting comments.
func PostCommentPrivateRoutes(route fiber.Router) {
	posts := route.Group("/posts/:post")

	posts.Post("/comments", middleware.JWTProtected(), controllers.AddComment)

	posts.Patch("/comments/:comment", middleware.JWTProtected(), controllers.UpdateComment)

	posts.Delete("/comments/:comment", middleware.JWTProtected(), controllers.DeleteComment)
}
