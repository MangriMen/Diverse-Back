package controllers

import (
	"fmt"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// swagger:route GET /posts Post getPosts
// Returns a list of all posts
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   200: GetPostsResponse
//   default: ErrorResponse

func GetPosts(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	posts, err := db.GetPosts()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "users not found",
			"count":   0,
			"users":   nil,
		})
	}

	postsToSend := lo.Map(posts, func(item models.DBPost, index int) models.Post {
		return helpers.PreparePostToSend(item, db)
	})

	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"count":   len(postsToSend),
		"posts":   postsToSend,
	})
}

// swagger:route GET /posts/{post} Post getPost
// Returns the post by given id
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   200: GetPostResponse
//   default: ErrorResponse

func GetPost(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("post"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	post, err := db.GetPost(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "users not found",
			"count":   0,
			"users":   nil,
		})
	}

	postToSend := helpers.PreparePostToSend(post, db)
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"post":    postToSend,
	})
}

// swagger:route POST /posts Post createPost
// Creates the post with given info
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   201: CreateUpdatePostResponse
//   default: ErrorResponse

func CreatePost(c *fiber.Ctx) error {
	post := &models.Post{}
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	newPost := &models.DBPost{
		BasePost: models.BasePost{
			Id:          uuid.New(),
			Content:     post.Content,
			Description: post.Description,
			Likes:       0,
			CreatedAt:   time.Now(),
		},
		UserId: post.User.Id,
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(newPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if err := db.CreatePost(newPost); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route PATCH /posts/{post} Post updatePost
// Update post by id with given fields
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   201: CreateUpdatePostResponse
//   default: ErrorResponse

func UpdatePost(c *fiber.Ctx) error {
	claims, err := helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postId, err := uuid.Parse(c.Params("post"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	post := &models.Post{}
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	foundPost, err := db.GetPost(postId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	if userId != foundPost.UserId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to update post",
		})
	}

	if foundPost.CreatedAt.Add(configs.PostEditTimeSinceCreated).UTC().Before(time.Now().UTC()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("can't edit post after %s", configs.PostEditTimeSinceCreated.String()),
		})
	}

	foundPost.Description = helpers.GetNotEmpty(post.Description, foundPost.Description)

	validate := helpers.NewValidator()
	if err := validate.Struct(foundPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.UpdatePost(foundPost.Id, &foundPost); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"erorr":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route DELETE /posts/{post} Post deletePost
// Delete post by id
//
// Schemes: http, https
//
// Produces:
//   - application/json
//
// Responses:
//   204: DeletePostResponse
//   default: ErrorResponse

func DeletePost(c *fiber.Ctx) error {
	claims, err := helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postIdToDelete, err := uuid.Parse(c.Params("post"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	post := &models.DBPost{BasePost: models.BasePost{Id: postIdToDelete}}
	validate := helpers.NewValidator()
	if err := validate.StructPartial(post, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"erorr":   true,
			"message": "book with this ID not found",
		})
	}

	foundPost, err := db.GetPost(postIdToDelete)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	if userId != foundPost.UserId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete post",
		})
	}

	if err := db.DeletePost(post.Id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// swagger:route POST /posts/{post}/comments Post addComment
// Add comment to the given post
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   201: CreateUpdateCommentResponse
//   default: ErrorResponse

func AddComment(c *fiber.Ctx) error {
	claims, err := helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postId, err := uuid.Parse(c.Params("post"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	comment := &models.Comment{}
	if err := c.BodyParser(comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if _, err = db.GetUser(userId); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this ID not found",
		})
	}

	if _, err := db.GetPost(postId); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	newComment := &models.DBComment{
		BaseComment: models.BaseComment{
			Id:        uuid.New(),
			Content:   comment.Content,
			CreatedAt: time.Now(),
		},
		PostId: postId,
		UserId: userId,
	}
	newComment.UpdatedAt = newComment.CreatedAt

	validate := helpers.NewValidator()
	if err := validate.Struct(newComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.AddComment(newComment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route PATCH /posts/{post}/comments/{comment} Post updateComment
// Update comment content by comment id with given post id
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   201: CreateUpdateCommentResponse
//   default: ErrorResponse

func UpdateComment(c *fiber.Ctx) error {
	claims, err := helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentId, err := uuid.Parse(c.Params("comment"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	comment := &models.Comment{}
	if err := c.BodyParser(comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	foundComment, err := db.GetComment(commentId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "comment with this ID not found",
		})
	}

	if userId != foundComment.UserId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to update comment",
		})
	}

	if foundComment.CreatedAt.Add(configs.PostCommentEditTimeSinceCreated).UTC().Before(time.Now().UTC()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("can't edit comment after %s", configs.PostCommentEditTimeSinceCreated.String()),
		})
	}

	foundComment.Content = helpers.GetNotEmpty(comment.Content, foundComment.Content)
	foundComment.UpdatedAt = time.Now()

	validate := helpers.NewValidator()
	if err := validate.Struct(foundComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.UpdateComment(foundComment.Id, &foundComment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"erorr":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route DELETE /posts/{post}/comments/{comment} Post deleteComment
// Delete comment by comment id with given post id
//
// Schemes: http, https
//
// Produces:
//   - application/json
//
// Responses:
//   204: DeleteCommentResponse
//   default: ErrorResponse

func DeleteComment(c *fiber.Ctx) error {
	claims, err := helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentIdToDelete, err := uuid.Parse(c.Params("comment"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	comment := &models.DBComment{BaseComment: models.BaseComment{Id: commentIdToDelete}}
	validate := helpers.NewValidator()
	if err := validate.StructPartial(comment, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"erorr":   true,
			"message": err.Error(),
		})
	}

	foundComment, err := db.GetComment(commentIdToDelete)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "comment with this ID not found",
		})
	}

	if userId != foundComment.UserId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete comment",
		})
	}

	if err := db.DeleteComment(comment.Id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
