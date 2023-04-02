package controllers_test

import (
	"fmt"
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

func TestAddRelation(t *testing.T) {
	tests := []struct {
		name         string
		desription   string
		route        string
		expectedCode int
	}{
		{
			name:         "Add relation success",
			desription:   "Get success response",
			route:        "/api/v1/users/id/relations",
			expectedCode: fiber.StatusCreated,
		},
		{
			name:         "Add relation error. Bad request",
			desription:   "Get success response",
			route:        "/api/v1/users/id/relations",
			expectedCode: fiber.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := server.InitAPI()

			_, token, err := helpers.RegisterUserForTest(app)
			if err != nil {
				log.Fatal(err.Error())
			}

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

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
