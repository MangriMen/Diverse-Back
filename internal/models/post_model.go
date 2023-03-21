package models

import (
	"time"

	"github.com/google/uuid"
)

type BasePost struct {
	// The id for this post
	// required: true
	Id uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

	// The url to the post content
	// required: true
	Content string `db:"content" json:"content" validate:"required"`

	// Post description
	// required: true
	Description string `db:"description" json:"description" validate:"lte=2048"`

	// Number of likes
	Likes int `db:"likes" json:"likes"`

	// The time the post was created
	// required: true
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DBPost struct {
	BasePost

	// The id of the user who created the post
	// required: true
	UserId uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
}

// Post represents the post for this application
// swagger:model
type Post struct {
	BasePost

	User *User `json:"user"`

	Comments []Comment `json:"comments"`
}

func (p *DBPost) ToPost() Post {
	return Post{BasePost: p.BasePost}
}

func (p *Post) ToDBPost() DBPost {
	return DBPost{BasePost: p.BasePost}
}
