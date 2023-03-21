package controllers

import (
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
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

	users, err := db.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "users not found",
			"count":   0,
			"users":   nil,
		})
	}

	users = lo.Map(users, func(item models.DBUser, index int) models.DBUser {
		item.PrepareToSend()
		return item
	})

	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"count":   len(users),
		"users":   users,
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
	userId, err := uuid.Parse(c.Params("user"))
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

	user, err := db.GetUser(userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with given id was not found",
			"user":    nil,
		})
	}

	user.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"user":    user,
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
	user := &models.DBUser{}
	if err := c.BodyParser(user); err != nil {
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

	foundUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email not found",
		})
	}

	if ok := helpers.CheckPasswordHash(user.Password, foundUser.Password); !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "wrong password",
		})
	}

	token, err := helpers.GenerateNewAccessToken(foundUser.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	foundUser.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"token":   token,
		"user":    foundUser,
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
	user := &models.DBUser{}
	if err := c.BodyParser(user); err != nil {
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

	if _, err := db.GetUserByEmail(user.Email); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email already exists",
		})
	}

	_, errEmail := db.GetUserByEmail(user.Email)
	_, errUsername := db.GetUserByUsername(user.Username)
	if errEmail == nil || errUsername == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email or username already exists",
		})
	}

	user.Id = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	validate := helpers.NewValidator()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "cannot create user",
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

	user.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"token":   token,
		"user":    user,
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

	user, err := db.GetUser(userId)
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

	user.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"token":   token,
		"user":    user,
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

	user := &models.DBUser{}
	if err := c.BodyParser(user); err != nil {
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

	foundUser, err := db.GetUser(userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this id not found",
		})
	}

	user.Id = foundUser.Id
	user.Email = helpers.GetNotEmpty(user.Email, foundUser.Email)
	user.Password = helpers.GetNotEmpty(user.Password, foundUser.Password)
	user.Username = helpers.GetNotEmpty(user.Username, foundUser.Username)
	user.Name = helpers.GetNotEmpty(user.Name, foundUser.Name)
	user.CreatedAt = foundUser.CreatedAt

	user.UpdatedAt = time.Now()

	validate := helpers.NewValidator()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	if err := db.UpdateUser(foundUser.Id, user); err != nil {
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

	userIdToDelete, err := uuid.Parse(c.Params("user"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if userId != userIdToDelete {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete user",
		})
	}

	user := &models.DBUser{BaseUser: models.BaseUser{Id: userIdToDelete}}
	validate := helpers.NewValidator()
	if err := validate.StructPartial(user, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"erorr":   true,
			"message": "book with this id not found",
		})
	}

	if err := db.DeleteUser(user.Id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
