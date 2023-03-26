// Package models provides base models for database and work with REST API
package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseComment represents a base comment struct in a system.
type BaseComment struct {
	// The id for this comment
	// required: true
	ID uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

	// Comment content
	// required: true
	Content string `db:"content" json:"content" validate:"required"`

	// The time the comment was created
	// required: true
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// The time the comment was updated
	// required: true
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// DBComment represents a comment struct from database.
type DBComment struct {
	BaseComment

	// Parent post id
	// required: true
	PostID uuid.UUID `db:"post_id" json:"post_id" validate:"required,uuid"`

	// Id of the user who wrote the comment
	// required: true
	UserID uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
}

// ToComment converts the DBComment to Comment model.
func (c *DBComment) ToComment() Comment {
	return Comment{BaseComment: c.BaseComment}
}

// Comment represents the comment for this application
// swagger:model
type Comment struct {
	BaseComment

	User *User `json:"user"`
}
