package controllers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/jwthelpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/userhelpers"
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
// Responses:
//   200: GetUsersResponse
//   default: ErrorResponse

// GetUsers is used to fetch users from database with request parameters.
func GetUsers(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
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
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(responses.GetUserResponseBody{
		User: dbUser.ToUser(),
	})
}

// swagger:route GET /users/username/{username} User getUserByUsername
// Returns the user by given username
//
// Responses:
//   200: GetUserResponse
//   default: ErrorResponse

// GetUserByUsername is used to fetch user from database by username.
func GetUserByUsername(c *fiber.Ctx) error {
	usernameIDParams, err := helpers.GetParamsAndValidate[parameters.UsernameIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbUser, err := db.GetUserByUsername(usernameIDParams.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(responses.GetUserResponseBody{
		User: dbUser.ToUser(),
	})
}

// swagger:route POST /login User loginUser
// Returns the user and token by given credentials
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
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundDBUser, err := db.GetUserByEmail(loginRequestBody.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	if ok := userhelpers.CheckPasswordHash(loginRequestBody.Password, foundDBUser.Password); !ok {
		return helpers.Response(c, fiber.StatusForbidden, configs.WrongEmailOrPasswordError)
	}

	token, err := jwthelpers.GenerateNewAccessToken(foundDBUser.ID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, "err.Error()")
	}

	return c.JSON(responses.RegisterLoginUserResponseBody{
		Token: token,
		User:  foundDBUser.ToUser(),
	})
}

// swagger:route POST /register User createUser
// Returns the user and token by given credentials
//
// Responses:
//   201: RegisterLoginUserResponse
//   default: ErrorResponse

// CreateUser is used to create new user.
func CreateUser(c *fiber.Ctx) error {
	registerRequestBody, err := helpers.GetBodyAndValidate[parameters.RegisterRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	_, errEmail := db.GetUserByEmail(registerRequestBody.Email)
	_, errUsername := db.GetUserByUsername(registerRequestBody.Username)
	if errEmail == nil || errUsername == nil {
		return helpers.Response(
			c,
			fiber.StatusConflict,
			configs.UserAlreadyExistsError,
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

	user.Password, err = userhelpers.HashPassword(registerRequestBody.Password)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(user); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.CreateUser(user); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	token, err := jwthelpers.GenerateNewAccessToken(user.ID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(
		responses.RegisterLoginUserResponseBody{
			Token: token,
			User:  user.ToUser(),
		})
}

// swagger:route Get /fetch User fetchUser
// Get user and new token if user exists
//
// Security:
//   bearerAuth:
//
// Responses:
//   200: RegisterLoginUserResponse
//   default: ErrorResponse

// FetchUser is used to regenerate user token and fetch info.
func FetchUser(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	dbUser, err := db.GetUser(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	token, err := jwthelpers.GenerateNewAccessToken(userID)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(responses.RegisterLoginUserResponseBody{
		Token: token,
		User:  dbUser.ToUser(),
	})
}

// swagger:route PATCH /users/{user} User updateUser
// Update user by id with given fields
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

	userIDParams, err := helpers.GetParamsAndValidate[parameters.UserIDParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	userUpdateRequestBody, err := helpers.GetBodyAndValidate[parameters.UserUpdateRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if userID != userIDParams.User {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundUser, err := db.GetUser(userIDParams.User)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundUser.Email = helpers.GetNotEmpty(userUpdateRequestBody.Email, foundUser.Email)
	foundUser.Username = helpers.GetNotEmpty(userUpdateRequestBody.Username, foundUser.Username)
	foundUser.Name = helpers.GetNotEmpty(userUpdateRequestBody.Name, foundUser.Name)
	foundUser.About = helpers.GetNotEmpty(userUpdateRequestBody.About, foundUser.About)

	foundUser.AvatarURL = helpers.GetNotEmpty(userUpdateRequestBody.AvatarURL, foundUser.AvatarURL)

	if userUpdateRequestBody.Password != "" {
		foundUser.Password, err = userhelpers.HashPassword(userUpdateRequestBody.Password)
		if err != nil {
			return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
		}
	}

	foundUser.UpdatedAt = time.Now()

	validate := helpers.NewValidator()
	if err = validate.Struct(foundUser); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.UpdateUser(&foundUser); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

// swagger:route PATCH /users/password User updateUserPassword
// Update user password
//
// Security:
//   bearerAuth:
//
// Responses:
//   201: UpdateUserPasswordResponse
//   default: ErrorResponse

// UpdateUserPassword is used to update user password by ID.
func UpdateUserPassword(c *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromToken(c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	userUpdateRequestBody, err := helpers.GetBodyAndValidate[parameters.UserUpdatePasswordRequestBody](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	foundUser, err := db.GetUser(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.Response(c, fiber.StatusNotFound, configs.UserNotFoundError)
		}

		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	if userhelpers.CheckPasswordHash(userUpdateRequestBody.OldPassword, foundUser.Password) {
		foundUser.Password, err = userhelpers.HashPassword(userUpdateRequestBody.Password)
		if err != nil {
			return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
		}
	} else {
		return helpers.Response(c, fiber.StatusForbidden, configs.WrongPassword)
	}

	foundUser.UpdatedAt = time.Now()

	validate := helpers.NewValidator()
	if err = validate.Struct(foundUser); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, helpers.ValidatorErrors(err))
	}

	if err = db.UpdateUser(&foundUser); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
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

	if userID != userIDParams.User {
		return helpers.Response(c, fiber.StatusForbidden, configs.ForbiddenError)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	if err = db.DeleteUser(userIDParams.User); err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
