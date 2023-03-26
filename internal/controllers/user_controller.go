package controllers

import (
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/jwthelpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// swagger:route GET /users User allowempty
// Returns a list of all users
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   200: GetUsersResponse
//   default: ErrorResponse

// GetUsers is used to fetch users from database with request parameters.
func GetUsers(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Users not found")
	}

	usersToSend := lo.Map(dbUsers, func(item models.DBUser, index int) models.User {
		return item.ToUser()
	})

	return c.JSON(responses.GetUsersResponseBody{
		Count: len(usersToSend),
		Users: usersToSend,
	})
}

// swagger:route GET /users/{user} User getUser
// Returns the user by given id
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   200: GetUserResponse
//   default: ErrorResponse

// GetUser is used to fetch user from database by ID.
func GetUser(c *fiber.Ctx) error {
	userIDParams, err := helpers.GetParamsAndValidate[parameters.UserIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbUser, err := db.GetUser(userIDParams.User)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "User with given id was not found")
	}

	return c.JSON(responses.GetUserResponseBody{
		User: dbUser.ToUser(),
	})
}

// swagger:route POST /login User loginUser
// Returns the user and token by given credentials
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   200: RegisterLoginUserResponse
//   default: ErrorResponse

// LoginUser is used to generate token to existing user.
func LoginUser(c *fiber.Ctx) error {
	loginRequestBody, err := helpers.GetBodyAndValidate[parameters.LoginRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	foundDBUser, err := db.GetUserByEmail(loginRequestBody.Email)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "User with this email not found")
	}

	if ok := helpers.CheckPasswordHash(loginRequestBody.Password, foundDBUser.Password); !ok {
		return helpers.Response(c, fiber.StatusForbidden, "Wrong email or password")
	}

	token, err := jwthelpers.GenerateNewAccessToken(foundDBUser.ID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(responses.RegisterLoginUserResponseBody{
		Token: token,
		User:  foundDBUser.ToUser(),
	})
}

// swagger:route POST /register User createUser
// Returns the user and token by given credentials
//
// Produces:
//   - application/json
//
// Schemes: http, https
//
// Responses:
//   200: RegisterLoginUserResponse
//   default: ErrorResponse

// CreateUser is used to create new user.
func CreateUser(c *fiber.Ctx) error {
	registerRequestBody, err := helpers.GetBodyAndValidate[parameters.RegisterRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	_, errEmail := db.GetUserByEmail(registerRequestBody.Email)
	_, errUsername := db.GetUserByUsername(registerRequestBody.Username)
	if errEmail == nil || errUsername == nil {
		return helpers.Response(
			c,
			fiber.StatusConflict,
			"User with this email or username already exists",
		)
	}

	user := &models.DBUser{
		BaseUser: models.BaseUser{
			ID:        uuid.New(),
			Email:     registerRequestBody.Email,
			Username:  registerRequestBody.Username,
			CreatedAt: time.Now(),
		},
	}
	user.UpdatedAt = user.CreatedAt

	user.Password, err = helpers.HashPassword(registerRequestBody.Password)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, "Cannot create user")
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(user); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.CreateUser(user); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	token, err := jwthelpers.GenerateNewAccessToken(user.ID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(
		responses.RegisterLoginUserResponseBody{
			Token: token,
			User:  user.ToUser(),
		})
}

// swagger:route Get /fetch User fetchUser
// Get user and new token if user exists
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
//   201: RegisterLoginUserResponse
//   default: ErrorResponse

// FetchUser is used to regenerate user token and fetch info.
func FetchUser(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	dbUser, err := db.GetUser(userID)
	if err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "User with this ID not found")
	}

	token, err := jwthelpers.GenerateNewAccessToken(userID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(responses.RegisterLoginUserResponseBody{
		Token: token,
		User:  dbUser.ToUser(),
	})
}

// swagger:route PATCH /users/{user} User updateUser
// Update user by id with given fields
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
//   201: UpdateUserResponse
//   default: ErrorResponse

// UpdateUser is used to update the user by ID.
func UpdateUser(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	userUpdateRequestBody, err := helpers.GetBodyAndValidate[parameters.UserUpdateRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	foundUser, err := db.GetUser(userID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, "User with this id not found")
	}

	foundUser.Email = helpers.GetNotEmpty(userUpdateRequestBody.Email, foundUser.Email)
	foundUser.Username = helpers.GetNotEmpty(userUpdateRequestBody.Username, foundUser.Username)
	foundUser.Name = helpers.GetNotEmpty(userUpdateRequestBody.Name, foundUser.Name)

	if userUpdateRequestBody.Password != "" {
		foundUser.Password, err = helpers.HashPassword(userUpdateRequestBody.Password)
		if err != nil {
			return helpers.Response(c, fiber.StatusInternalServerError, "Cannot update password")
		}

		foundUser.UpdatedAt = time.Now()

		validate := helpers.NewValidator()
		if err = validate.Struct(foundUser); err != nil {
			return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
		}

		if err = db.UpdateUser(&foundUser); err != nil {
			return helpers.Response(c, fiber.StatusInternalServerError, err)
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route DELETE /users/{user} User deleteUser
// Delete user by id
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
//   204: DeleteUserResponse
//   default: ErrorResponse

// DeleteUser is used to delete the user by ID.
func DeleteUser(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	userIDParams, err := helpers.GetParamsAndValidate[parameters.UserIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(userIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if userID != userIDParams.User {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete user",
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"erorr":   true,
			"message": "book with this id not found",
		})
	}

	if err = db.DeleteUser(userIDParams.User); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
