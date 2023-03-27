package parameters

import "github.com/google/uuid"

// UserIDParams includes the id of the user.
type UserIDParams struct {
	// in: path
	// required: true
	User uuid.UUID `params:"user" json:"user"`
}

// UserIDRequest is used to represent a request that requires a user id parameter,
// such as fetching a specific user, updating user, or deleting a user.
// swagger:parameters getUser updateUser deleteUser
type UserIDRequest struct {
	UserIDParams
}

// LoginRequestBody includes the email and password of the user.
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

// LoginRequest is used for login a user.
// swagger:parameters loginUser
type LoginRequest struct {
	// in: body
	// required: true
	Body LoginRequestBody
}

// RegisterRequestBody includes the email, username, and password of the user being created.
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

// RegisterRequest is used for register a new user.
// swagger:parameters createUser
type RegisterRequest struct {
	// in: body
	// required: true
	Body RegisterRequestBody
}

// UserUpdateRequestBody includes the new email, username, name or password for the user.
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

// UserUpdateRequest represents a request to update a user's information,
// including fields such as email, username, name or password
// swagger:parameters updateUser
type UserUpdateRequest struct {
	// in: body
	// required: true
	Body UserUpdateRequestBody
}
