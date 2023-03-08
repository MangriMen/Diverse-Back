package models

import (
	"github.com/google/uuid"
)

// swagger: model
type DefaultResponseBody struct {
	// required: true
	Error bool `json:"error"`
	// required: true
	Message string `json:"message"`
}

// swagger:response
type ErrorResponse struct {
	// in: body
	Body DefaultResponseBody
}

// swagger:response
type GetUsersResponse struct {
	// in: body
	Body struct {
		DefaultResponseBody
		// required: true
		Users []BaseUser `json:"users"`
		// required: true
		Count int `json:"count"`
	}
}

// swagger:response
type GetUserResponse struct {
	// in: body
	Body struct {
		DefaultResponseBody
		// required: true
		User BaseUser `json:"user"`
	}
}

// swagger:response
type RegisterLoginUserResponse struct {
	// in: body
	Body struct {
		DefaultResponseBody
		// required: true
		User BaseUser `json:"user"`
		// required: true
		Token string `json:"token"`
	}
}

// swagger:response
type UpdateUserResponse string

// swagger:response
type DeleteUserResponse struct {
}

// swagger:parameters getUser updateUser deleteUser
type UserIdParameter struct {
	// in: path
	// required: true
	Id uuid.UUID `json:"id"`
}

// swagger:parameters loginUser
type LoginParameters struct {
	// in: body
	Body struct {
		// required: true
		// min length: 6
		// max length: 255
		Email string `json:"email"`
		// required: true
		// min length: 8
		// max length: 256
		Password string `json:"password"`
	}
}

// swagger:parameters createUser
type RegisterParameters struct {
	// in: body
	// required: true
	Body struct {
		// required: true
		// min length: 6
		// max length: 255
		Email string `json:"email"`
		// required: true
		// min length: 8
		// max length: 256
		Password string `json:"password"`
		// required: true
		// min length: 1
		// max length: 32
		Username string `json:"username"`
	}
}
