package controllers

import (
	"time"

	"github.com/MangriMen/Value-Back/api/database"
	"github.com/MangriMen/Value-Back/internal/helpers"
	"github.com/MangriMen/Value-Back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

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

func GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
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

	user, err := db.GetUser(id)
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

func LoginUser(c *fiber.Ctx) error {
	user := &models.User{}
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

	token, err := helpers.GenerateNewAccessToken()
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

func CreateUser(c *fiber.Ctx) error {
	user := &models.User{}
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

	token, err := helpers.GenerateNewAccessToken()
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

func UpdateUser(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := helpers.ExtractTokenMetadata(c)
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

	user := &models.User{}
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

	foundUser, err := db.GetUser(user.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user with this ID not found",
		})
	}

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

func DeleteUser(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := helpers.ExtractTokenMetadata(c)
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

	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

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
			"message": "book with this ID not found",
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
