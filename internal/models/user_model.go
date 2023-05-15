package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseUser represents a base user struct in a system.
type BaseUser struct {
	// The id for this user
	// required: true
	ID uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

	// The email for this user
	// required: true
	Email string `db:"email" json:"email" validate:"required,gte=6,lte=255"`

	// The username for this user
	// required: true
	Username string `db:"username" json:"username" validate:"required,gte=1,lte=32"`

	// The name for this user
	Name string `db:"name" json:"name" validate:"lte=32"`

	// The time the user was registered
	// required: true
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// The time the user was updated
	// required: true
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// URL to user avatar
	AvatarURL *string `db:"avatar_url" json:"avatar_url"`

	About *string `db:"about" json:"about"`
}

// DBUser represents a user struct from database.
type DBUser struct {
	BaseUser

	// The password for this user
	// required: true
	Password string `db:"password" json:"password,omitempty" validate:"required,gte=8,lte=256"`
}

// ToUser converts the DBUser to User model.
func (u *DBUser) ToUser() User {
	return User{BaseUser: u.BaseUser}
}

// User represents the user for this application
// swagger:model
type User struct {
	BaseUser
}
