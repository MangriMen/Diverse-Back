// Package controllers_test provides test for app controllers.
package controllers_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name               string
		route              string
		expectedStatusCode int
	}{
		{
			name:               "Get users",
			route:              "/api/v1/users/",
			expectedStatusCode: 200,
		},
	}

	if !helpers.IsRunningInContainer() {
		helpers.LoadEnvironment(".env")
	}

	app := server.InitAPI()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, tt.route, nil)

			resp, _ := app.Test(req, -1)

			defer helpers.CloseQuietly(resp.Body)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			bodyString := string(bodyBytes)

			assert.Equalf(t, tt.expectedStatusCode, resp.StatusCode, bodyString)
		})
	}
}
