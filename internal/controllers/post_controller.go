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
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetPostsResponse
//   default: ErrorResponse

// GetPosts is used to fetch posts from database with request parameters.
func GetPosts(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

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

	filter, err := posthelpers.GenerateFilter(userID, postsFetchRequestQuery, db)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbPosts, err := db.GetPosts(postsFetchRequestQuery, filter)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postsToSend := lo.Map(dbPosts, func(item models.DBPost, index int) models.Post {
		return posthelpers.PreparePostToSend(item, userID, db)
	})

	return c.JSON(responses.GetPostsResponseBody{
		Count: len(postsToSend),
		Data:  postsToSend,
	})
}

// swagger:route GET /posts/{post} Post getPost
// Returns the post by given ID
//
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetPostResponse
//   default: ErrorResponse

// GetPost is used to fetch post from database by ID.
func GetPost(c *fiber.Ctx) error {
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

	dbPost, err := db.GetPost(postIDParams.Post)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.PostNotFoundError)
	}

	postToSend := posthelpers.PreparePostToSend(dbPost, userID, db)

	return c.JSON(responses.GetPostResponseBody{
		Data: postToSend,
	})
}

// swagger:route POST /posts Post createPost
// Creates the post with given info
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
// Security:
//   bearerAuth:
//
// Responses:
//   201: GetPostResponse
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
		return helpers.Response(c, fiber.StatusNotFound, configs.PostNotFoundError)
	}

	if userID != foundPost.UserID {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	if foundPost.CreatedAt.Add(configs.PostEditTimeSinceCreated).UTC().
		Before(time.Now().UTC()) {
		return helpers.Response(c, fiber.StatusForbidden, fmt.Sprintf(
			configs.CantEditAfterErrorFormat,
			"post",
			configs.PostEditTimeSinceCreated.String(),
		))
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
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postToSend := posthelpers.PreparePostToSend(foundPost, userID, db)

	return c.JSON(responses.GetPostResponseBody{
		Data: postToSend,
	})
}

// swagger:route POST /posts/{post}/like Post likePost
// Set like to the post by ID
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

	like := &models.DBPostLike{
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
		return helpers.Response(c, fiber.StatusNotFound, configs.PostNotFoundError)
	}

	postToSend := posthelpers.PreparePostToSend(foundPost, userID, db)

	return c.JSON(responses.GetPostResponseBody{
		Data: postToSend,
	})
}

// swagger:route DELETE /posts/{post}/like Post unlikePost
// Unset like to the post by ID
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

	like := &models.DBPostLike{
		PostID: postIDParams.Post,
		UserID: userID,
	}

	if err = db.UnlikePost(like); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundPost, err := db.GetPost(like.PostID)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.PostNotFoundError)
	}

	postToSend := posthelpers.PreparePostToSend(foundPost, userID, db)

	return c.JSON(responses.GetPostResponseBody{
		Data: postToSend,
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
		return helpers.Response(c, fiber.StatusNotFound, configs.PostNotFoundError)
	}

	if userID != foundPost.UserID {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	if err = db.DeletePost(foundPost.ID); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
