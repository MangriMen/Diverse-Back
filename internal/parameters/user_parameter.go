package parameters

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
	// required: true
	Body struct {
		// required: true
		// min length: 6
		// max length: 255
		Email string `json:"email" validate:"required,gte=6,lte=255"`

		// required: true
		// min length: 8
		// max length: 256
		Password string `json:"password" validate:"required,gte=8,lte=256"`
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
}

// swagger:parameters updateUser
type UserUpdateParameters struct {
	// in: body
	// required: true
	Body struct {
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
}
