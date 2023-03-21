package swagger_models

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
)

type BaseResponseBody struct {
	// required: true
	Error bool `json:"error"`
	// required: true
	Message string `json:"message"`
}

// swagger:response
type ErrorResponse struct {
	// in: body
	Body BaseResponseBody
}

// swagger:response
type GetUsersResponse struct {
	// in: body
	Body struct {
		BaseResponseBody
		// required: true
		Users []models.User `json:"users"`
		// required: true
		Count int `json:"count"`
	}
}

// swagger:response
type GetUserResponse struct {
	// in: body
	Body struct {
		BaseResponseBody
		// required: true
		User models.User `json:"user"`
	}
}

// swagger:response
type RegisterLoginUserResponse struct {
	// in: body
	Body struct {
		BaseResponseBody
		// required: true
		User models.User `json:"user"`
		// required: true
		Token string `json:"token"`
	}
}

// swagger:response
type UpdateUserResponse string

// swagger:response
type DeleteUserResponse struct {
}
