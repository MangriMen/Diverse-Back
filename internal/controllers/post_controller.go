// Package controllers provide functionality for web application routes
package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/posthelpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
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

// GetPosts is used to fetch posts from database with request parameters.
func GetPosts(c *fiber.Ctx) error {
	postsFetchRequestQuery, err := helpers.GetQueryAndValidate[parameters.PostsFetchRequestQuery](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if postsFetchRequestQuery.LastSeenPostCreatedAt.IsZero() {
		postsFetchRequestQuery.LastSeenPostCreatedAt = time.Now()
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbPosts, err := db.GetPosts(postsFetchRequestQuery)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Posts not found")
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
// Returns the post by given ID
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

// GetPost is used to fetch post from database by ID.
func GetPost(c *fiber.Ctx) error {
	postIDParams, err := helpers.GetParamsAndValidate[parameters.PostIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbPost, err := db.GetPost(postIDParams.Post)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Post with given id was not found")
	}

	postToSend := posthelpers.PreparePostToSend(dbPost, db)

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

// CreatePost is used to create a new post.
func CreatePost(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postCreateRequestBody, err := helpers.GetBodyAndValidate[parameters.PostCreateRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	newPost := &models.DBPost{
		BasePost: models.BasePost{
			ID:          uuid.New(),
			Content:     postCreateRequestBody.Content,
			Description: postCreateRequestBody.Description,
			Likes:       0,
			CreatedAt:   time.Now(),
		},
		UserID: userID,
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(newPost); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	if err = db.CreatePost(newPost); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route PATCH /posts/{post} Post updatePost
// Update post by ID with given fields
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

// UpdatePost is used to update the post by ID.
func UpdatePost(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postIDParams, err := helpers.GetParamsAndValidate[parameters.PostIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	postUpdateRequestBody, err := helpers.GetBodyAndValidate[parameters.PostUpdateRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
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

	validate := helpers.NewValidator()
	if err = validate.Struct(foundPost); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.UpdatePost(&foundPost); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"erorr":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route POST /posts/{post}/like Post likePost
// Set like to the post by ID
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
//   201: GetPostResponse
//   default: ErrorResponse

// LikePost is used to like the post by ID.
func LikePost(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postIDParams, err := helpers.GetParamsAndValidate[parameters.PostIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	like := &models.DBLike{
		ID:     uuid.New(),
		PostID: postIDParams.Post,
		UserID: userID,
	}

	if err = db.LikePost(like); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == configs.DBDuplicateError {
			return helpers.Response(c, fiber.StatusConflict, err.Error())
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundPost, err := db.GetPost(like.PostID)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Post with this ID not found")
	}

	postToSend := posthelpers.PreparePostToSend(foundPost, db)

	return c.JSON(responses.GetPostResponseBody{
		Post: postToSend,
	})
}

// swagger:route DELETE /posts/{post}/like Post unlikePost
// Set like to the post by ID
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
//   201: GetPostResponse
//   default: ErrorResponse

// UnlikePost is used to unlike the post by ID.
func UnlikePost(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postIDParams, err := helpers.GetParamsAndValidate[parameters.PostIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	like := &models.DBLike{
		PostID: postIDParams.Post,
		UserID: userID,
	}

	if err = db.UnlikePost(like); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundPost, err := db.GetPost(like.PostID)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Post with this ID not found")
	}

	postToSend := posthelpers.PreparePostToSend(foundPost, db)

	return c.JSON(responses.GetPostResponseBody{
		Post: postToSend,
	})
}

// swagger:route DELETE /posts/{post} Post deletePost
// Delete post by ID
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

// DeletePost is used to delete the post by ID.
func DeletePost(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postIDParams, err := helpers.GetParamsAndValidate[parameters.PostIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	foundPost, err := db.GetPost(postIDParams.Post)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Post with this ID not found")
	}

	if userID != foundPost.UserID {
		return helpers.Response(c, fiber.StatusNotFound, "Not enough permission to delete post")
	}

	if err = db.DeletePost(foundPost.ID); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
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

// GetComments is used to fetch the post comments by post ID.
func GetComments(c *fiber.Ctx) error {
	postCommentIDParams, err := helpers.GetParamsAndValidate[parameters.PostCommentIDParams](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	commentsFetchRequestQuery, err := helpers.GetQueryAndValidate[parameters.CommentsFetchRequestQuery](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if commentsFetchRequestQuery.LastSeenCommentCreatedAt.IsZero() {
		commentsFetchRequestQuery.LastSeenCommentCreatedAt = time.Now()
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbPosts, err := db.GetComments(postCommentIDParams.Post, commentsFetchRequestQuery)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Posts not found")
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

// AddComment is used to add the comment to the post by ID.
func AddComment(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	commentAddRequestParams, err := helpers.GetParamsAndValidate[parameters.CommentAddRequestParams](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	commentAddRequestBody, err := helpers.GetBodyAndValidate[parameters.CommentAddRequestBody](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	if _, err = db.GetUser(userID); err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "User with this ID not found")
	}

	if _, err = db.GetPost(commentAddRequestParams.Post); err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Post with this ID not found")
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

	validate := helpers.NewValidator()
	if err = validate.Struct(newComment); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.AddComment(newComment); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route PATCH /posts/{post}/comments/{comment} Post updateComment
// Update comment content by comment ID with given post ID
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

// UpdateComment is used to update the comment on the post by post ID and comment ID.
func UpdateComment(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postCommentIDParams, err := helpers.GetParamsAndValidate[parameters.PostCommentIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	commentUpdateRequestBody, err := helpers.GetBodyAndValidate[parameters.CommentUpdateRequestBody](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundComment, err := db.GetComment(postCommentIDParams.Comment)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Comment with this ID not found")
	}

	if userID != foundComment.UserID {
		return helpers.Response(
			c,
			fiber.StatusNotFound,
			"Not enough permission to update comment",
		)
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

	validate := helpers.NewValidator()
	if err = validate.Struct(foundComment); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.UpdateComment(&foundComment); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route DELETE /posts/{post}/comments/{comment} Post deleteComment
// Delete comment by comment ID with given post ID
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

// DeleteComment is used to delete the comment on the post by post ID and comment ID.
func DeleteComment(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postCommentIDParams, err := helpers.GetParamsAndValidate[parameters.PostCommentIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	foundComment, err := db.GetComment(postCommentIDParams.Comment)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Comment with this ID not found")
	}

	if userID != foundComment.UserID {
		return helpers.Response(
			c,
			fiber.StatusNotFound,
			"Not enough permission to delete comment",
		)
	}

	if err = db.DeleteComment(foundComment.ID); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
