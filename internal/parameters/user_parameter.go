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

// UsernameIDParams includes the id of the user.
type UsernameIDParams struct {
	// in: path
	// required: true
	Username string `params:"username" json:"username"`
}

// UsernameIDRequest is used to represent a request that requires a user id parameter,
// such as fetching a specific user, updating user, or deleting a user.
// swagger:parameters getUserByUsername
type UsernameIDRequest struct {
	UsernameIDParams
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
	Email string `json:"email" validate:"omitempty,gte=6,lte=255"`

	// min length: 1
	// max length: 32
	Username string `json:"username" validate:"omitempty,gte=1,lte=32"`

	// max length: 32
	Name string `json:"name" validate:"omitempty,lte=32"`

	// min length: 8
	// max length: 256
	Password string `json:"password" validate:"omitempty,gte=8,lte=256"`

	AvatarURL *string `json:"avatar_url" validate:"omitempty"`

	// min length: 0
	// max length: 2048
	About *string `db:"about" json:"about"`
}

// UserUpdateRequest represents a request to update a user's information,
// including fields such as email, username, name or password
// swagger:parameters updateUser
type UserUpdateRequest struct {
	UserIDParams
	// in: body
	// required: true
	Body UserUpdateRequestBody
}

// UserUpdatePasswordRequestBody includes the old and new password of the user.
type UserUpdatePasswordRequestBody struct {
	// min length: 8
	// max length: 256
	OldPassword string `json:"old_password" validate:"required,gte=8,lte=256"`

	// min length: 8
	// max length: 256
	Password string `json:"password" validate:"required,gte=8,lte=256"`
}

// UserUpdatePasswordRequest represents a request to update a user's password.
// swagger:parameters updateUserPassword
type UserUpdatePasswordRequest struct {
	// in: body
	// required: true
	Body UserUpdatePasswordRequestBody
}
