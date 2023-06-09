package responses

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
)

// GetUsersResponseBody includes the slice of users.
type GetUsersResponseBody struct {
	BaseResponseBody
	// required: true
	Count int `json:"count"`
	// required: true
	Users []models.User `json:"users"`
}

// GetUsersResponse represent the response retrived on get users request.
// swagger:response
type GetUsersResponse struct {
	// in: body
	Body GetUsersResponseBody
}

// GetUserResponseBody includes the signle user for a given ID.
type GetUserResponseBody struct {
	BaseResponseBody
	// required: true
	User models.User `json:"user"`
}

// GetUserResponse represent the response retrived on get user request.
// swagger:response
type GetUserResponse struct {
	// in: body
	Body GetUserResponseBody
}

// RegisterLoginUserResponseBody includes the token and user model when user register or login.
type RegisterLoginUserResponseBody struct {
	BaseResponseBody
	// required: true
	Token string `json:"token"`
	// required: true
	User models.User `json:"user"`
}

// RegisterLoginUserResponse represent the response retrived on register or login request.
// swagger:response
type RegisterLoginUserResponse struct {
	// in: body
	Body RegisterLoginUserResponseBody
}

// UpdateUserResponse represents response for successfully update user request.
// swagger:response
type UpdateUserResponse string

// UpdateUserPasswordResponse represents response for successfully update user password.
// swagger:response
type UpdateUserPasswordResponse string

// DeleteUserResponse represents response for successfully delete user request.
// swagger:response
type DeleteUserResponse struct {
}
