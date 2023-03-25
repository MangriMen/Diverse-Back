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

func GetUsers(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	dbUsers, err := db.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "users not found",
			"count":   0,
			"users":   nil,
		})
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

func GetUser(c *fiber.Ctx) error {
	userIDParams := &parameters.UserIDParams{}
	if err := c.ParamsParser(userIDParams); err != nil {
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

	dbUser, err := db.GetUser(userIDParams.User)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with given id was not found",
			"user":    nil,
		})
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

func LoginUser(c *fiber.Ctx) error {
	loginRequestBody := &parameters.LoginRequestBody{}
	if err := c.BodyParser(loginRequestBody); err != nil {
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

	foundDBUser, err := db.GetUserByEmail(loginRequestBody.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email not found",
		})
	}

	if ok := helpers.CheckPasswordHash(loginRequestBody.Password, foundDBUser.Password); !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "wrong password",
		})
	}

	token, err := jwthelpers.GenerateNewAccessToken(foundDBUser.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
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

func CreateUser(c *fiber.Ctx) error {
	registerRequestBody := &parameters.RegisterRequestBody{}
	if err := c.BodyParser(registerRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(registerRequestBody); err != nil {
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

	_, errEmail := db.GetUserByEmail(registerRequestBody.Email)
	_, errUsername := db.GetUserByUsername(registerRequestBody.Username)
	if errEmail == nil || errUsername == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email or username already exists",
		})
	}

	user := &models.DBUser{
		BaseUser: models.BaseUser{
			Id:        uuid.New(),
			Email:     registerRequestBody.Email,
			Username:  registerRequestBody.Username,
			CreatedAt: time.Now(),
		},
	}
	user.UpdatedAt = user.CreatedAt

	user.Password, err = helpers.HashPassword(registerRequestBody.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "cannot create user",
		})
	}

	if err = validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	token, err := jwthelpers.GenerateNewAccessToken(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
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

func FetchUser(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if claims.Expires < now {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unathorized, check expiration time of your token",
		})
	}

	userID, err := uuid.Parse(claims.ID)
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

	dbUser, err := db.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this id not found",
		})
	}

	token, err := jwthelpers.GenerateNewAccessToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
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

func UpdateUser(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if claims.Expires < now {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unathorized, check expiration time of your token",
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userUpdateRequestBody := &parameters.UserUpdateRequestBody{}
	if err = c.BodyParser(userUpdateRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err = validate.Struct(userUpdateRequestBody); err != nil {
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

	foundUser, err := db.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this id not found",
		})
	}

	foundUser.Email = helpers.GetNotEmpty(userUpdateRequestBody.Email, foundUser.Email)
	foundUser.Username = helpers.GetNotEmpty(userUpdateRequestBody.Username, foundUser.Username)
	foundUser.Name = helpers.GetNotEmpty(userUpdateRequestBody.Name, foundUser.Name)

	if userUpdateRequestBody.Password != "" {
		foundUser.Password, err = helpers.HashPassword(userUpdateRequestBody.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "cannot update password",
			})
		}
	}

	foundUser.UpdatedAt = time.Now()

	if err = validate.Struct(foundUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err = db.UpdateUser(&foundUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"erorr":   true,
			"message": err.Error(),
		})
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

func DeleteUser(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if claims.Expires < now {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unathorized, check expiration time of your token",
		})
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userIDParams := &parameters.UserIDParams{}
	if err = c.ParamsParser(userIDParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
