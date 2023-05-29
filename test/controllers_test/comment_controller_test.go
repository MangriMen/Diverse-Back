// Package controllers_test provides test for app controllers.
package controllers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetCommentsCount(t *testing.T) {
	tests := []struct {
		name               string
		route              string
		expectedStatusCode int
	}{
		{
			name:               "Get comments count without authorization",
			route:              "/posts/4202cc54-779a-4508-8af6-4b4dda99bca5/comments/count",
			expectedStatusCode: 200,
		},
	}

	app := server.InitAPI()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, tt.route, nil)

			resp, _ := app.Test(req, 1)

			assert.Equalf(t, tt.expectedStatusCode, resp.StatusCode, resp.Status)
		})
	}
}
