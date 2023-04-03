package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddRelation(t *testing.T) {
	tests := []struct {
		name         string
		relationType models.RelationType
		desription   string
		route        string
		expectedCode int
	}{
		{
			name:         "Add follower relation success",
			relationType: models.Follower,
			desription:   "Get success response",
			route:        "/api/v1/users/%s/relations",
			expectedCode: fiber.StatusCreated,
		},
		{
			name:         "Add following relation success",
			relationType: models.Following,
			desription:   "Get success response",
			route:        "/api/v1/users/%s/relations",
			expectedCode: fiber.StatusCreated,
		},
		{
			name:         "Add blocked relation success",
			relationType: models.Blocked,
			desription:   "Get success response",
			route:        "/api/v1/users/%s/relations",
			expectedCode: fiber.StatusCreated,
		},
		{
			name:         "Add relation error. Bad request",
			relationType: "failed",
			desription:   "Get success response",
			route:        "/api/v1/users/%s/relations",
			expectedCode: fiber.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := server.InitAPI()

			user, token, err := helpers.RegisterUserForTest(app)
			if err != nil {
				log.Fatal(err.Error())
			}

			relationUser, _, err := helpers.RegisterUserForTest(app)
			if err != nil {
				log.Fatal(err.Error())
			}

			relationAddRequestBody := &parameters.RelationAddRequestBody{
				RelationUserID: relationUser.ID,
				Type:           tt.relationType,
			}

			var buf bytes.Buffer
			if err = json.NewEncoder(&buf).Encode(relationAddRequestBody); err != nil {
				log.Fatal(err.Error())
			}

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf(tt.route, user.ID),
				&buf,
			)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, configs.TestResponseTimeout)
			if err != nil {
				log.Fatal(err.Error())
			}

			rawBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err.Error())
			}

			body, err := helpers.ParseResponseBody[responses.BaseResponseBody](rawBody)
			if err != nil {
				log.Fatal(err.Error())
			}

			message := helpers.GetMessageFromResponseBody(body, tt.desription)

			assert.Equalf(t, tt.expectedCode, resp.StatusCode, message)
		})
	}
}
