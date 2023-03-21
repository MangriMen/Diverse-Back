package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseComment struct {
	// The id for this comment
	// required: true
	Id uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

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

type DBComment struct {
	BaseComment

	// Parent post id
	// required: true
	PostId uuid.UUID `db:"post_id" json:"post_id" validate:"required,uuid"`

	// Id of the user who wrote the comment
	// required: true
	UserId uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
}

type Comment struct {
	BaseComment

	User *User `json:"user"`
}

func (c *DBComment) ToComment() Comment {
	return Comment{BaseComment: c.BaseComment}
}
