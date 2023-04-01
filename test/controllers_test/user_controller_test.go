package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
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
			expectedCode: fiber.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helpers.LoadEnvironment(".env")

			app := server.InitAPI()

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, 100)

			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.desription)
		})
	}
}
