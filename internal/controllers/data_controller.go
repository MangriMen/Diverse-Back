package controllers

import (
	"os"
	"path/filepath"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/datahelpers"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
)

// swagger:route POST /data Data uploadData
// Returns the ID of uploaded data
//
// Responses:
//   200: UploadDataResponse
//   default: ErrorResponse

// UploadData is used to upload data files.
func UploadData(c *fiber.Ctx) error {
	receivedFile, err := c.FormFile("file")
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	baseType, _, err := datahelpers.ValidateContentType(
		receivedFile,
		configs.GetAllowedMIMEBaseTypes(),
	)
	if err != nil {
		return helpers.Response(c, fiber.StatusUnsupportedMediaType, err.Error())
	}

	filename := datahelpers.GenerateUniqueFilename()
	pathToFolder := filepath.Join(configs.DataPath, baseType)

	err = os.MkdirAll(pathToFolder, os.ModePerm)
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	err = datahelpers.ProcessFile(receivedFile, filepath.Join(pathToFolder, filename))
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(
		responses.UploadDataResponseBody{
			Path: filepath.Join("/data", baseType, "raw", filename),
		},
	)
}
