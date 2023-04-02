package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
)

// ParseResponseBody is used to unmarshall response body to T.
// Returns T object.
func ParseResponseBody[T responses.ResponseBody](
	resp *http.Response,
) (T, error) {
	var response T

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// GetMessageFromResponseBody is used to get message from BaseResponseBody or
// if nil returns default message.
func GetMessageFromResponseBody(
	body responses.BaseResponseBody,
	defaultMessage string,
) string {
	if body.Message != nil {
		if responseMessage, ok := body.Message.(string); ok {
			return responseMessage
		}
	}

	return defaultMessage
}

// RegisterUserForTest is used to register temp user for using token with private routes.
func RegisterUserForTest(app *fiber.App) (*models.User, string, error) {
	registerRequestBody := &parameters.RegisterRequestBody{
		Email:    "test@gmail.com",
		Username: "test",
		Password: "testtest",
	}

	stringBody, err := json.Marshal(&registerRequestBody)
	if err != nil {
		return nil, "", err
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(stringBody))

	resp, err := app.Test(req, 500)
	if err != nil {
		return nil, "", err
	}

	body, err := ParseResponseBody[responses.BaseResponseBody](resp)
	if err != nil {
		return nil, "", err
	}

	if body.Error && body.Message != nil {
		return nil, "", fmt.Errorf(body.Message.(string))
	}

	trueBody, err := ParseResponseBody[responses.RegisterLoginUserResponseBody](resp)
	if err != nil {
		return nil, "", err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, "", err
	}

	return &trueBody.User, trueBody.Token, nil
}
