package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Email     string    `db:"email" json:"email" validate:"required,lte=255"`
	Password  string    `db:"password" json:"password,omitempty" validate:"required,lte=256"`
	Username  string    `db:"username" json:"username" validate:"required,lte=32"`
	Name      string    `db:"name" json:"name" validate:"lte=32"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (user *User) PrepareToSend() {
	user.Password = ""
}
