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

// swagger:route GET /posts/{post}/comments Post getComments
// Returns a list of post comments
//
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetCommentsResponse
//   default: ErrorResponse

// GetComments is used to fetch the post comments by post ID.
func GetComments(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

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
		return helpers.Response(c, fiber.StatusNotFound, configs.CommentsNotFoundError)
	}

	commentsToSend := lo.Map(dbPosts, func(item models.DBComment, index int) models.Comment {
		return posthelpers.PrepareCommentToPost(item, userID, db)
	})

	return c.JSON(responses.GetCommentsResponseBody{
		Count: len(commentsToSend),
		Data:  commentsToSend,
	})
}

// swagger:route POST /posts/{post}/comments Post addComment
// Add comment to the given post
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
		return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
	}

	if _, err = db.GetPost(commentAddRequestParams.Post); err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.PostNotFoundError)
	}

	newComment := &models.DBComment{
		BaseComment: models.BaseComment{
			ID:        uuid.New(),
			Content:   commentAddRequestBody.Content,
			CreatedAt: time.Now(),
			Likes:     0,
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
		return helpers.Response(c, fiber.StatusNotFound, configs.CommentNotFoundError)
	}

	if userID != foundComment.UserID {
		return helpers.Response(
			c,
			fiber.StatusForbidden,
			configs.ForbiddenError,
		)
	}

	if foundComment.CreatedAt.Add(configs.PostCommentEditTimeSinceCreated).UTC().
		Before(time.Now().UTC()) {
		return helpers.Response(c, fiber.StatusForbidden, fmt.Sprintf(
			configs.CantEditAfterErrorFormat,
			"comment",
			configs.PostCommentEditTimeSinceCreated.String(),
		))
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

	commentToSend := posthelpers.PrepareCommentToPost(foundComment, userID, db)

	return c.JSON(responses.GetCommentResponseBody{
		Data: commentToSend,
	})
}

// swagger:route POST /posts/{post}/comments/{comment}/like Post likeComment
// Set like to the comment by ID
//
// Security:
//   bearerAuth:
//
// Responses:
//   201: GetCommentResponse
//   default: ErrorResponse

// LikeComment is used to like the comment by ID.
func LikeComment(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postCommentIDParams, err := helpers.GetParamsAndValidate[parameters.PostCommentIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	like := &models.DBCommentLike{
		ID:        uuid.New(),
		CommentID: postCommentIDParams.Comment,
		UserID:    userID,
	}

	if err = db.LikeComment(like); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == configs.DBDuplicateError {
			return helpers.Response(c, fiber.StatusConflict, err.Error())
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundComment, err := db.GetComment(like.CommentID)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.CommentNotFoundError)
	}

	commentToSend := posthelpers.PrepareCommentToPost(foundComment, userID, db)

	return c.JSON(responses.GetCommentResponseBody{
		Data: commentToSend,
	})
}

// swagger:route DELETE /posts/{post}/comments/{comment}/like Post unlikeComment
// Unset like to the comment by ID
//
// Security:
//   bearerAuth:
//
// Responses:
//   201: GetCommentResponse
//   default: ErrorResponse

// UnlikeComment is used to unlike the comment by ID.
func UnlikeComment(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	postCommentIDParams, err := helpers.GetParamsAndValidate[parameters.PostCommentIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	like := &models.DBCommentLike{
		CommentID: postCommentIDParams.Comment,
		UserID:    userID,
	}

	if err = db.UnlikeComment(like); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundComment, err := db.GetComment(like.CommentID)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.CommentNotFoundError)
	}

	commentToSend := posthelpers.PrepareCommentToPost(foundComment, userID, db)

	return c.JSON(responses.GetCommentResponseBody{
		Data: commentToSend,
	})
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
		return helpers.Response(c, fiber.StatusNotFound, configs.CommentNotFoundError)
	}

	if userID != foundComment.UserID {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	if err = db.DeleteComment(foundComment.ID); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
