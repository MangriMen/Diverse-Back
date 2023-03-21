package responses

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
)

type GetUsersResponseBody struct {
	BaseResponseBody
	// required: true
	Count int `json:"count"`
	// required: true
	Users []models.User `json:"users"`
}

// swagger:response
type GetUsersResponse struct {
	// in: body
	Body GetUsersResponseBody
}

type GetUserResponseBody struct {
	BaseResponseBody
	// required: true
	User models.User `json:"user"`
}

// swagger:response
type GetUserResponse struct {
	// in: body
	Body GetUserResponseBody
}

type RegisterLoginUserResponseBody struct {
	BaseResponseBody
	// required: true
	Token string `json:"token"`
	// required: true
	User models.User `json:"user"`
}

// swagger:response
type RegisterLoginUserResponse struct {
	// in: body
	Body RegisterLoginUserResponseBody
}

// swagger:response
type UpdateUserResponse string

// swagger:response
type DeleteUserResponse struct {
}
