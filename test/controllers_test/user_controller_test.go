package controllers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name         string
		desription   string
		route        string
		expectedCode int
	}{
		{
			name:         "Get users success",
			desription:   "Get success response",
			route:        "/api/v1/users",
			expectedCode: fiber.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helpers.LoadEnvironment(".env")

			app := server.InitAPI()

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, 500)
			if err != nil {
				log.Fatal(err.Error())
			}

			body, err := helpers.ParseResponseBody[responses.BaseResponseBody](resp)
			if err != nil {
				log.Fatal(err.Error())
			}

			message := helpers.GetMessageFromResponseBody(body, tt.desription)

			assert.Equalf(t, tt.expectedCode, resp.StatusCode, message)
		})
	}
}
