// Package controllers_test provides test for app controllers.
package controllers_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/jwthelpers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCommentsCount(t *testing.T) {
	tests := []struct {
		name               string
		route              string
		userID             string
		expectedStatusCode int
	}{
		// {
		// 	name:               "Get comments count without authorization",
		// 	route:              "/api/v1/posts/023753e6-6539-4c95-93e6-1d469d80f28d/comments/count/",
		// 	expectedStatusCode: 403,
		// },
		{
			name:               "Get comments count",
			route:              "/api/v1/posts/023753e6-6539-4c95-93e6-1d469d80f28d/comments/count/",
			userID:             "17afe1a4-aaed-4263-a29b-781389509cb6",
			expectedStatusCode: 200,
		},
	}

	if !helpers.IsRunningInContainer() {
		helpers.LoadEnvironment(".env")
	}

	app := server.InitAPI()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			var err error
			if tt.userID != "" {
				token, err = jwthelpers.GenerateNewAccessToken(uuid.MustParse(tt.userID))
				if err != nil {
					panic(err)
				}
			} else {
				token = ""
			}

			req := httptest.NewRequest(fiber.MethodGet, tt.route, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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
