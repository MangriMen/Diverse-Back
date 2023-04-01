package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
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

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)

			resp, _ := app.Test(req, 1)

			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.desription)
		})
	}
}
