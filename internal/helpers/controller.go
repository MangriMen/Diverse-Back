// Package helpers provides a set of helper functions for handling various tasks
package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/MangriMen/Diverse-Back/internal/helpers/jwthelpers"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUserIDFromToken parses the JWT token and returns the user id.
func GetUserIDFromToken(c *fiber.Ctx) (uuid.UUID, error) {
	userID := uuid.UUID{}

	claims, err := jwthelpers.GetTokenClaims(c)
	if err != nil {
		return userID, err
	}

	userID, err = uuid.Parse(claims.ID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

// GetParamsAndValidate parses parameters from the request to the structure, validate it and returns.
func GetParamsAndValidate[T parameters.RequestParams](
	c *fiber.Ctx,
) (*T, error) {
	var params T
	if err := c.ParamsParser(&params); err != nil {
		return nil, err
	}

	validate := NewValidator()
	if err := validate.Struct(&params); err != nil {
		errorMessage, marshalErr := json.Marshal(ValidatorErrors(err))
		if marshalErr != nil {
			return nil, marshalErr
		}
		return nil, fmt.Errorf(string(errorMessage))
	}

	return &params, nil
}

// GetQueryAndValidate parses query from the request to the structure, validate it and returns.
func GetQueryAndValidate[T parameters.RequestQuery](
	c *fiber.Ctx,
) (*T, error) {
	var query T
	if err := c.QueryParser(&query); err != nil {
		return nil, err
	}

	validate := NewValidator()
	if err := validate.Struct(&query); err != nil {
		errorMessage, marshalErr := json.Marshal(ValidatorErrors(err))
		if marshalErr != nil {
			return nil, marshalErr
		}
		return nil, fmt.Errorf(string(errorMessage))
	}

	return &query, nil
}

// GetBodyAndValidate parses body from the request to the structure, validate it and returns.
func GetBodyAndValidate[T parameters.RequestBody](
	c *fiber.Ctx,
) (*T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return nil, err
	}

	validate := NewValidator()
	if err := validate.Struct(&body); err != nil {
		errorMessage, marshalErr := json.Marshal(ValidatorErrors(err))
		if marshalErr != nil {
			return nil, marshalErr
		}
		return nil, fmt.Errorf(string(errorMessage))
	}

	return &body, nil
}

// Response used to compact return default error response.
func Response(c *fiber.Ctx, responseStatus int, err interface{}) error {
	return c.Status(responseStatus).JSON(responses.BaseResponseBody{
		Error:   true,
		Message: &err,
	})
}
