package controllers

import (
	"errors"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/userhelpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/samber/lo"
)

// swagger:route GET /users/{user}/relations User getRelations
// Returns a list of users from given relation
//
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetRelationsResponse
//   default: ErrorResponse

// GetRelations is used to fetch relation between users from database with request parameters.
func GetRelations(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	userIDParams, err := helpers.GetParamsAndValidate[parameters.UserIDParams](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationGetRequestQuery, err := helpers.GetQueryAndValidate[parameters.RelationGetRequestQuery](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if relationGetRequestQuery.LastSeenRelationCreatedAt.IsZero() {
		relationGetRequestQuery.LastSeenRelationCreatedAt = time.Now()
	}

	if userID != userIDParams.User {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbRelations, err := db.GetRelations(userIDParams.User, relationGetRequestQuery)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	relationsToSend := lo.Map(
		dbRelations,
		func(item models.DBRelation, index int) models.Relation {
			return *userhelpers.PrepareRelationToSend(item, db)
		},
	)

	return c.JSON(responses.GetRelationResponseBody{
		Count:     len(relationsToSend),
		Relations: relationsToSend,
	})
}

// swagger:route GET /users/{user}/relations/{relationUser} User getRelationStatus
// Returns a true or false from given relation
//
// Security:
//   bearerAuth:
//
// Responses:
//   200: GetRelationStatusResponse
//   default: ErrorResponse

// GetRelationStatus is used to fetch relation existence between users from database with request parameters.
func GetRelationStatus(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationGetStatusParams, err := helpers.GetParamsAndValidate[parameters.RelationGetStatusParams](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if userID != relationGetStatusParams.User {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	relationStatus, err := db.GetRelationStatus(relationGetStatusParams)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	statusToSend := userhelpers.PrepareRelationStatusToSend(relationStatus)

	return c.JSON(responses.GetRelationStatusResponseBody{
		Follower:  statusToSend[models.Follower],
		Following: statusToSend[models.Following],
		Blocked:   statusToSend[models.Blocked],
	})
}

// swagger:route POST /users/{user}/relations/{relationUser} User addRelation
// Add realtion with given info
//
// Security:
//   bearerAuth:
//
// Responses:
//   201: AddRelationResponse
//   default: ErrorResponse

// AddRelation is used to add relation between users with request parameters.
func AddRelation(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationGetStatusParams, err := helpers.GetParamsAndValidate[parameters.RelationGetStatusParams](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationAddDeleteRequestBody, err := helpers.GetBodyAndValidate[parameters.RelationAddDeleteRequestBody](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if userID != relationGetStatusParams.User {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	relation := &models.DBRelation{
		BaseRelation: models.BaseRelation{
			ID:        uuid.New(),
			Type:      relationAddDeleteRequestBody.Type,
			CreatedAt: time.Now(),
		},
		UserID:         relationGetStatusParams.User,
		RelationUserID: relationGetStatusParams.RelationUser,
	}

	if err = db.AddRelation(relation); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == configs.DBDuplicateError {
			return helpers.Response(c, fiber.StatusConflict, err.Error())
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route DELETE /users/{user}/relations/{relationUser} User deleteRelation
// Returns relation by id
//
// Security:
//   bearerAuth:
//
// Responses:
//   204: DeleteRelationResponse
//   default: ErrorResponse

// DeleteRelation is used to delete relation between users.
func DeleteRelation(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationGetStatusParams, err := helpers.GetParamsAndValidate[parameters.RelationGetStatusParams](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationAddDeleteRequestBody, err := helpers.GetBodyAndValidate[parameters.RelationAddDeleteRequestBody](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if userID != relationGetStatusParams.User {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	if err = db.DeleteRelation(relationGetStatusParams, relationAddDeleteRequestBody); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
