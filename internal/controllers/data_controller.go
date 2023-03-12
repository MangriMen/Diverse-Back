package controllers

import (
	"os"
	"path/filepath"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/gofiber/fiber/v2"
)

func UploadData(c *fiber.Ctx) error {
	receivedFile, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	baseType, _, err := helpers.ValidateContentType(receivedFile, configs.UploadMIMEBaseTypes)
	if err != nil {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	filename := helpers.GenerateUniqueFilename()
	url := filepath.Join(baseType, filename)
	pathToFolder := filepath.Join(configs.DATA_PATH, baseType)

	err = os.MkdirAll(pathToFolder, os.ModePerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	err = c.SaveFile(receivedFile, filepath.Join(pathToFolder, filename))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.GetResponse(err, helpers.DefaultError))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": nil,
		"url":     url,
	})
}
