package models

import (
	"time"

	"github.com/google/uuid"
)

// BasePost represents a base post struct in a system.
type BasePost struct {
	// The id for this post
	// required: true
	ID uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

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

// DBPost represents a post struct from database.
type DBPost struct {
	BasePost

	// The id of the user who created the post
	// required: true
	UserID uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
}

// ToPost converts the DBPost to Post model.
func (p *DBPost) ToPost() Post {
	return Post{BasePost: p.BasePost}
}

// Post represents the post for this application
// swagger:model
type Post struct {
	BasePost

	User *User `json:"user"`

	Comments []Comment `json:"comments"`
}

// DBLike represents a like struct from database.
type DBLike struct {
	// The id for this like
	// required: true
	ID uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

	// Parent post id
	// required: true
	PostID uuid.UUID `db:"post_id" json:"post_id" validate:"required,uuid"`

	// Id of the user who wrote the comment
	// required: true
	UserID uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
}

// Like represents the like for this application
// swagger:model
type Like struct {
	User User
}
