package controllers

import (
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/jwt_helpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// swagger:route GET /users getUsers
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
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	users, err := db.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "users not found",
		})
	}

	users = lo.Map(users, func(item models.User, index int) models.User {
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

// swagger:route GET /users/{id} getUser
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
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	user, err := db.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with given id was not found",
		})
	}

	user.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"user":    user,
	})
}

// swagger:route POST /login loginUser
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
	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
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

	token, err := jwt_helpers.GenerateNewAccessToken(foundUser.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	foundUser.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"token":   token,
		"user":    foundUser,
	})
}

// swagger:route POST /register createUser
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
	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	if _, err := db.GetUserByEmail(user.Email); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "user with this email already exists",
		})
	}

	if _, err := db.GetUserByUsername(user.Username); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "user with this username already exists",
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
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	if err := db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	token, err := jwt_helpers.GenerateNewAccessToken(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	user.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"token":   token,
		"user":    user,
	})
}

// swagger:route Get /fetch fetchUser
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
	claims, err := jwt_helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	user, err := db.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this ID not found",
		})
	}

	token, err := jwt_helpers.GenerateNewAccessToken(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	user.PrepareToSend()
	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"token":   token,
		"user":    user,
	})
}

// swagger:route PATCH /users/{id} updateUser
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
	claims, err := jwt_helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	foundUser, err := db.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this ID not found",
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

// swagger:route DELETE /users/{id} deleteUser
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
	claims, err := jwt_helpers.GetTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	idToDelete, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	if id != idToDelete {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "not enough permission to delete user",
		})
	}

	user := &models.User{BaseUser: models.BaseUser{Id: idToDelete}}
	validate := helpers.NewValidator()
	if err := validate.StructPartial(user, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": helpers.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	if err := db.DeleteUser(user.Id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
