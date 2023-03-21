package controllers

import (
	"fmt"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
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
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetPostsResponse
//   default: ErrorResponse

func GetPosts(c *fiber.Ctx) error {
	fetchPostParameters := &parameters.PostsFetchParameters{}
	if err := c.QueryParser(fetchPostParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(fetchPostParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if fetchPostParameters.LastSeenPostCreatedAt.IsZero() {
		fetchPostParameters.LastSeenPostCreatedAt = time.Now()
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	dbPosts, err := db.GetPosts(fetchPostParameters)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "posts not found",
		})
	}

	postsToSend := lo.Map(dbPosts, func(item models.DBPost, index int) models.Post {
		return helpers.PreparePostToSend(item, db)
	})

	return c.JSON(responses.GetPostsResponseBody{
		Count: len(postsToSend),
		Posts: postsToSend,
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
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetPostResponse
//   default: ErrorResponse

func GetPost(c *fiber.Ctx) error {
	postIdParameter := &parameters.PostIdParameter{}
	if err := c.QueryParser(postIdParameter); err != nil {
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

	post, err := db.GetPost(postIdParameter.Post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post not found",
		})
	}

	postToSend := helpers.PreparePostToSend(post, db)

	return c.JSON(responses.GetPostResponseBody{
		Post: postToSend,
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
// Security:
//   bearerAuth:
//
// Responses:
//   201: CreateUpdatePostResponse
//   default: ErrorResponse

func CreatePost(c *fiber.Ctx) error {
	postCreateParameters := &parameters.PostCreateParameters{}
	if err := c.QueryParser(postCreateParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postCreateParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	newPost := &models.DBPost{
		BasePost: models.BasePost{
			Id:          uuid.New(),
			Content:     postCreateParameters.Body.Content,
			Description: postCreateParameters.Body.Description,
			Likes:       0,
			CreatedAt:   time.Now(),
		},
		UserId: postCreateParameters.Body.UserId,
	}

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
// Security:
//   bearerAuth:
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

	postUpdateParameters := &parameters.PostUpdateParameters{}
	if err := c.QueryParser(postUpdateParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postUpdateParameters); err != nil {
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

	foundPost, err := db.GetPost(postUpdateParameters.Post)
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

	foundPost.Description = helpers.GetNotEmpty(postUpdateParameters.Body.Description, foundPost.Description)

	if err := validate.Struct(foundPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.UpdatePost(&foundPost); err != nil {
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
// Security:
//   bearerAuth:
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

	postIdParameter := &parameters.PostIdParameter{}
	if err := c.QueryParser(postIdParameter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postIdParameter); err != nil {
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

	foundPost, err := db.GetPost(postIdParameter.Post)
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

	if err := db.DeletePost(foundPost.Id); err != nil {
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
// Security:
//   bearerAuth:
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

	commentAddParametersParams := &parameters.CommentAddParametersParams{}
	if err := c.ParamsParser(commentAddParametersParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentAddParametersBody := &parameters.CommentAddParametersBody{}
	if err := c.BodyParser(commentAddParametersBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(commentAddParametersParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := validate.Struct(commentAddParametersBody); err != nil {
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

	if _, err = db.GetUser(userId); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this ID not found",
		})
	}

	if _, err := db.GetPost(commentAddParametersParams.Post); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	newComment := &models.DBComment{
		BaseComment: models.BaseComment{
			Id:        uuid.New(),
			Content:   commentAddParametersBody.Content,
			CreatedAt: time.Now(),
		},
		PostId: commentAddParametersParams.Post,
		UserId: userId,
	}
	newComment.UpdatedAt = newComment.CreatedAt

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
// Security:
//   bearerAuth:
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

	commentUpdateParameters := &parameters.CommentUpdateParameters{}
	if err := c.QueryParser(commentUpdateParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(commentUpdateParameters); err != nil {
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

	foundComment, err := db.GetComment(commentUpdateParameters.Comment)
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

	foundComment.Content = helpers.GetNotEmpty(commentUpdateParameters.Body.Content, foundComment.Content)
	foundComment.UpdatedAt = time.Now()

	if err := validate.Struct(foundComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.UpdateComment(&foundComment); err != nil {
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
// Security:
//   bearerAuth:
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

	postCommentIdParameter := &parameters.PostCommentIdParameter{}
	if err := c.QueryParser(postCommentIdParameter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postCommentIdParameter); err != nil {
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

	foundComment, err := db.GetComment(postCommentIdParameter.Comment)
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

	if err := db.DeleteComment(foundComment.Id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
