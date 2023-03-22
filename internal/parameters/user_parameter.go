package parameters

import "github.com/google/uuid"

type UserIdParams struct {
	// in: path
	// required: true
	User uuid.UUID `params:"user" json:"user"`
}

// swagger:parameters getUser updateUser deleteUser
type UserIdRequest struct {
	UserIdParams
}

type LoginRequestBody struct {
	// required: true
	// min length: 6
	// max length: 255
	Email string `json:"email" validate:"required,gte=6,lte=255"`

	// required: true
	// min length: 8
	// max length: 256
	Password string `json:"password" validate:"required,gte=8,lte=256"`
}

// swagger:parameters loginUser
type LoginRequest struct {
	// in: body
	// required: true
	Body LoginRequestBody
}

type RegisterRequestBody struct {
	// required: true
	// min length: 6
	// max length: 255
	Email string `json:"email" validate:"required,gte=6,lte=255"`

	// required: true
	// min length: 1
	// max length: 32
	Username string `json:"username" validate:"required,gte=1,lte=32"`

	// required: true
	// min length: 8
	// max length: 256
	Password string `json:"password" validate:"required,gte=8,lte=256"`
}

// swagger:parameters createUser
type RegisterRequest struct {
	// in: body
	// required: true
	Body RegisterRequestBody
}

type UserUpdateRequestBody struct {
	// min length: 6
	// max length: 255
	Email string `json:"email" validate:"gte=6,lte=255"`

	// min length: 1
	// max length: 32
	Username string `json:"username" validate:"gte=1,lte=32"`

	// max length: 32
	Name string `json:"name" validate:"lte=32"`

	// min length: 8
	// max length: 256
	Password string `json:"password" validate:"gte=8,lte=256"`
}

// swagger:parameters updateUser
type UserUpdateRequest struct {
	// in: body
	// required: true
	Body UserUpdateRequestBody
}
