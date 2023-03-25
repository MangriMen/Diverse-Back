// Package controllers provide functionality for web application routes
package controllers

import (
	"fmt"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/jwthelpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/posthelpers"
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
	postsFetchRequestQuery := &parameters.PostsFetchRequestQuery{}
	if err := c.QueryParser(postsFetchRequestQuery); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postsFetchRequestQuery); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if postsFetchRequestQuery.LastSeenPostCreatedAt.IsZero() {
		postsFetchRequestQuery.LastSeenPostCreatedAt = time.Now()
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	dbPosts, err := db.GetPosts(postsFetchRequestQuery)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "posts not found",
		})
	}

	postsToSend := lo.Map(dbPosts, func(item models.DBPost, index int) models.Post {
		return posthelpers.PreparePostToSend(item, db)
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
	postIDParams := &parameters.PostIDParams{}
	if err := c.QueryParser(postIDParams); err != nil {
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

	post, err := db.GetPost(postIDParams.Post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post not found",
		})
	}

	postToSend := posthelpers.PreparePostToSend(post, db)

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
	postCreateRequestBody := &parameters.PostCreateRequestBody{}
	if err := c.QueryParser(postCreateRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postCreateRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	newPost := &models.DBPost{
		BasePost: models.BasePost{
			ID:          uuid.New(),
			Content:     postCreateRequestBody.Content,
			Description: postCreateRequestBody.Description,
			Likes:       0,
			CreatedAt:   time.Now(),
		},
		UserID: postCreateRequestBody.UserID,
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

	if err = db.CreatePost(newPost); err != nil {
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
	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postIDParams := &parameters.PostIDParams{}
	if err = c.ParamsParser(postIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postUpdateRequestBody := &parameters.PostUpdateRequestBody{}
	if err = c.QueryParser(postUpdateRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(postIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = validate.Struct(postUpdateRequestBody); err != nil {
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

	foundPost, err := db.GetPost(postIDParams.Post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	if userID != foundPost.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to update post",
		})
	}

	if foundPost.CreatedAt.Add(configs.PostEditTimeSinceCreated).UTC().Before(time.Now().UTC()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"message": fmt.Sprintf(
				"can't edit post after %s",
				configs.PostEditTimeSinceCreated.String(),
			),
		})
	}

	foundPost.Description = helpers.GetNotEmpty(
		postUpdateRequestBody.Description,
		foundPost.Description,
	)

	if err = validate.Struct(foundPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = db.UpdatePost(&foundPost); err != nil {
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
	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postIDParams := &parameters.PostIDParams{}
	if err = c.QueryParser(postIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(postIDParams); err != nil {
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

	foundPost, err := db.GetPost(postIDParams.Post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	if userID != foundPost.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete post",
		})
	}

	if err = db.DeletePost(foundPost.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// swagger:route GET /posts/{post}/comments Post getComments
// Returns a list of post comments
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

func GetComments(c *fiber.Ctx) error {
	postCommentIDParams := &parameters.PostCommentIDParams{}
	if err := c.QueryParser(postCommentIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentsFetchRequestQuery := &parameters.CommentsFetchRequestQuery{}
	if err := c.QueryParser(commentsFetchRequestQuery); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(postCommentIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := validate.Struct(commentsFetchRequestQuery); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if commentsFetchRequestQuery.LastSeenCommentCreatedAt.IsZero() {
		commentsFetchRequestQuery.LastSeenCommentCreatedAt = time.Now()
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	dbPosts, err := db.GetComments(postCommentIDParams.Post, commentsFetchRequestQuery)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "posts not found",
		})
	}

	commentsToSend := lo.Map(dbPosts, func(item models.DBComment, index int) models.Comment {
		return posthelpers.PrepareCommentToPost(item, db)
	})

	return c.JSON(responses.GetCommentsResponseBody{
		Count:    len(commentsToSend),
		Comments: commentsToSend,
	})
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
	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentAddRequestParams := &parameters.CommentAddRequestParams{}
	if err = c.ParamsParser(commentAddRequestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentAddRequestBody := &parameters.CommentAddRequestBody{}
	if err = c.BodyParser(commentAddRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(commentAddRequestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = validate.Struct(commentAddRequestBody); err != nil {
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

	if _, err = db.GetUser(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this ID not found",
		})
	}

	if _, err := db.GetPost(commentAddRequestParams.Post); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "post with this ID not found",
		})
	}

	newComment := &models.DBComment{
		BaseComment: models.BaseComment{
			ID:        uuid.New(),
			Content:   commentAddRequestBody.Content,
			CreatedAt: time.Now(),
		},
		PostID: commentAddRequestParams.Post,
		UserID: userID,
	}
	newComment.UpdatedAt = newComment.CreatedAt

	if err = validate.Struct(newComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = db.AddComment(newComment); err != nil {
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
	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postCommentIDParams := &parameters.PostCommentIDParams{}
	if err = c.ParamsParser(postCommentIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	commentUpdateRequestBody := &parameters.CommentUpdateRequestBody{}
	if err = c.BodyParser(commentUpdateRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(postCommentIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = validate.Struct(commentUpdateRequestBody); err != nil {
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

	foundComment, err := db.GetComment(postCommentIDParams.Comment)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "comment with this ID not found",
		})
	}

	if userID != foundComment.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to update comment",
		})
	}

	if foundComment.CreatedAt.Add(configs.PostCommentEditTimeSinceCreated).
		UTC().
		Before(time.Now().UTC()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"message": fmt.Sprintf(
				"can't edit comment after %s",
				configs.PostCommentEditTimeSinceCreated.String(),
			),
		})
	}

	foundComment.Content = helpers.GetNotEmpty(
		commentUpdateRequestBody.Content,
		foundComment.Content,
	)
	foundComment.UpdatedAt = time.Now()

	if err = validate.Struct(foundComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = db.UpdateComment(&foundComment); err != nil {
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
	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	postCommentIDParams := &parameters.PostCommentIDParams{}
	if err = c.ParamsParser(postCommentIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(postCommentIDParams); err != nil {
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

	foundComment, err := db.GetComment(postCommentIDParams.Comment)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "comment with this ID not found",
		})
	}

	if userID != foundComment.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete comment",
		})
	}

	if err = db.DeleteComment(foundComment.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
