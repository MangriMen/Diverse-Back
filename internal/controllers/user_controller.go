package controllers

import (
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
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
	userIdParameter := &parameters.UserIdParameter{}
	if err := c.QueryParser(userIdParameter); err != nil {
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

	dbUser, err := db.GetUser(userIdParameter.User)
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
	loginParameters := &parameters.LoginParameters{}
	if err := c.BodyParser(loginParameters); err != nil {
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

	foundDBUser, err := db.GetUserByEmail(loginParameters.Body.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email not found",
		})
	}

	if ok := helpers.CheckPasswordHash(loginParameters.Body.Password, foundDBUser.Password); !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "wrong password",
		})
	}

	token, err := helpers.GenerateNewAccessToken(foundDBUser.Id)
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
	registerParameters := &parameters.RegisterParameters{}
	if err := c.BodyParser(registerParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(registerParameters); err != nil {
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

	_, errEmail := db.GetUserByEmail(registerParameters.Body.Email)
	_, errUsername := db.GetUserByUsername(registerParameters.Body.Username)
	if errEmail == nil || errUsername == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email or username already exists",
		})
	}

	user := &models.DBUser{
		BaseUser: models.BaseUser{
			Id:        uuid.New(),
			Email:     registerParameters.Body.Email,
			Username:  registerParameters.Body.Username,
			CreatedAt: time.Now(),
		},
	}
	user.UpdatedAt = user.CreatedAt

	user.Password, err = helpers.HashPassword(registerParameters.Body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "cannot create user",
		})
	}

	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	token, err := helpers.GenerateNewAccessToken(user.Id)
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

	claims, err := helpers.GetTokenMetadata(c)
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

	userId, err := uuid.Parse(claims.Id)
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

	dbUser, err := db.GetUser(userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this id not found",
		})
	}

	token, err := helpers.GenerateNewAccessToken(userId)
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

	claims, err := helpers.GetTokenMetadata(c)
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

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userUpdateParameters := &parameters.UserUpdateParameters{}
	if err := c.BodyParser(userUpdateParameters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(userUpdateParameters); err != nil {
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

	foundUser, err := db.GetUser(userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this id not found",
		})
	}

	foundUser.Email = helpers.GetNotEmpty(userUpdateParameters.Body.Email, foundUser.Email)
	foundUser.Username = helpers.GetNotEmpty(userUpdateParameters.Body.Username, foundUser.Username)
	foundUser.Name = helpers.GetNotEmpty(userUpdateParameters.Body.Name, foundUser.Name)

	if userUpdateParameters.Body.Password != "" {
		foundUser.Password, err = helpers.HashPassword(userUpdateParameters.Body.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "cannot update password",
			})
		}
	}

	foundUser.UpdatedAt = time.Now()

	if err := validate.Struct(foundUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.UpdateUser(&foundUser); err != nil {
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

	claims, err := helpers.GetTokenMetadata(c)
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

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	userIdParameter := &parameters.UserIdParameter{}
	if err := c.QueryParser(userIdParameter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := helpers.NewValidator()
	if err := validate.Struct(userIdParameter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if userId != userIdParameter.User {
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

	if err := db.DeleteUser(userIdParameter.User); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
