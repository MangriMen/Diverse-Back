package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/MangriMen/Diverse-Back/internal/helpers/jwthelpers"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUserIDFromToken parse JWT token and returns user ID.
func GetUserIDFromToken(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := jwthelpers.GetTokenMetadata(c)
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(claims.ID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

// GetParamsAndValidate get parameters from request to struct, validate it and returns.
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

// GetQueryAndValidate get query from request to struct, validate it and returns.
func GetQueryAndValidate[T parameters.RequestQuery](
	c *fiber.Ctx,
) (*T, error) {
	var params T
	if err := c.QueryParser(&params); err != nil {
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

// GetBodyAndValidate get body from request to struct, validate it and returns.
func GetBodyAndValidate[T parameters.RequestBody](
	c *fiber.Ctx,
) (*T, error) {
	var params T
	if err := c.BodyParser(&params); err != nil {
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
