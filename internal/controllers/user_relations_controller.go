package controllers

import (
	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/userhelpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// swagger:route GET /users/{user}/relations User getRelations
// Returns a list of users from given relation
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
//   200: GetRelationResponse
//   default: ErrorResponse

// GetRelations is used to fetch relation between users from database with request parameters.
func GetRelations(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	relationGetRequestQuery, err := helpers.GetQueryAndValidate[parameters.RelationGetRequestQuery](
		c,
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbRelations, err := db.GetRelations(userID, relationGetRequestQuery)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.RelationsNotFoundError)
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

// swagger:route GET /users/{user}/relations/{relationUser} User getRelationByUser
// Returns a true or false from given relation
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
//   200: GetRelationByUserResponse
//   default: ErrorResponse

// IsRelationWithUser is used to fetch relation existence between users from database with request parameters.
func IsRelationWithUser(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.UsersNotFoundError)
	}

	usersToSend := lo.Map(dbUsers, func(item models.DBUser, index int) models.User {
		return item.ToUser()
	})

	return c.JSON(responses.GetUsersResponseBody{
		Count: len(usersToSend),
		Users: usersToSend,
	})
}

// swagger:route POST /users/{user}/relations User addRelation
// Add realtion with given info
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
//   200: AddRelationResponse
//   default: ErrorResponse

// AddRelation is used to add relation between users with request parameters.
func AddRelation(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.UsersNotFoundError)
	}

	usersToSend := lo.Map(dbUsers, func(item models.DBUser, index int) models.User {
		return item.ToUser()
	})

	return c.JSON(responses.GetUsersResponseBody{
		Count: len(usersToSend),
		Users: usersToSend,
	})
}

// swagger:route DELETE /users/{user}/relations/{relation} User deleteRelation
// Returns relation by id
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
//   200: DeleteRelationResponse
//   default: ErrorResponse

// DeleteRelation is used to delete relation between users.
func DeleteRelation(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, configs.UsersNotFoundError)
	}

	usersToSend := lo.Map(dbUsers, func(item models.DBUser, index int) models.User {
		return item.ToUser()
	})

	return c.JSON(responses.GetUsersResponseBody{
		Count: len(usersToSend),
		Users: usersToSend,
	})
}
