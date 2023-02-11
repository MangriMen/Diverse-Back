package controllers

import (
	"github.com/MangriMen/Value-Back/api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	return c.JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"user":    user,
	})
}
