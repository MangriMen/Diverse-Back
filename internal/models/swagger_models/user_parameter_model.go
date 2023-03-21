package swagger_models

import "github.com/google/uuid"

// swagger:parameters getUser updateUser deleteUser
type UserIdParameter struct {
	// in: path
	// required: true
	User uuid.UUID `json:"user"`
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
