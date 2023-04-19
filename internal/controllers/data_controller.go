package controllers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/datahelpers"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
)

// swagger:route GET /data/{type}/{image} Data getData
// Returns the requested data
//
// Responses:
//   304: SuccessResponse
//   default: ErrorResponse

// GetData is used to get data with parameters.
func GetData(c *fiber.Ctx) error {
	getDataRequestParams, err := helpers.GetParamsAndValidate[parameters.GetDataRequestParams](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	getDataRequestQuery, err := helpers.GetQueryAndValidate[parameters.GetDataRequestQuery](c)
	if err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error())
	}

	if getDataRequestParams.Type == configs.MIMEBaseImage {
		file, fileErr := bimg.Read(
			filepath.Join(
				configs.DataPath,
				getDataRequestParams.Type,
				getDataRequestParams.Image,
			),
		)
		if fileErr != nil {
			return fileErr
		}

		imageMetadata, metadataErr := bimg.Metadata(file)
		if metadataErr != nil {
			return metadataErr
		}

		var options bimg.Options

		switch {
		case getDataRequestQuery.Width != nil &&
			*getDataRequestQuery.Width < imageMetadata.Size.Width:
			options = bimg.Options{Width: *getDataRequestQuery.Width}
		case getDataRequestQuery.Height != nil &&
			*getDataRequestQuery.Height < imageMetadata.Size.Height:
			options = bimg.Options{Height: *getDataRequestQuery.Height}
		}

		image, imageErr := bimg.NewImage(file).Process(options)
		if imageErr != nil {
			return imageErr
		}

		c.Set(
			"Content-Type",
			strings.Join(
				[]string{
					getDataRequestParams.Type,
					bimg.ImageTypeName(bimg.DetermineImageType(image)),
				},
				"/",
			),
		)

		return c.Send(image)
	}

	return helpers.Response(c, fiber.StatusBadRequest, "Invalid type")
}

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
			Path: filepath.Join("/data", baseType, filename),
		},
	)
}
