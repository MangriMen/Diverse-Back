package parameters

import (
	"time"

	"github.com/google/uuid"
)

// PostIDParams includes the id of the post.
type PostIDParams struct {
	// in: path
	// required: true
	Post uuid.UUID `params:"post" json:"post" validate:"required"`
}

// PostIDRequest is used to represent a request that requires a post id parameter,
// such as fetching a specific post or deleting a post.
// swagger:parameters getPost updatePost deletePost
type PostIDRequest struct {
	PostIDParams
}

// PostCreateRequestBody includes the content, description and id of the user creating the post.
type PostCreateRequestBody struct {
	// required: true
	Content string `json:"content" validate:"required"`

	// required: true
	// max length: 2048
	Description string `json:"description" validate:"lte=2048"`

	// required: true
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

// PostCreateRequest is used for creating a new post.
// It includes the new post content.
// swagger:parameters createPost
type PostCreateRequest struct {
	// in: body
	// required: true
	Body PostCreateRequestBody
}

// PostUpdateRequestBody includes the new description for the post.
type PostUpdateRequestBody struct {
	// required: true
	// max length: 2048
	Description string `json:"description" validate:"lte=2048"`
}

// PostUpdateRequest is used for updating an existing post.
// It includes the ID of the post to be updated as well as the new post content.
type PostUpdateRequest struct {
	PostIDParams

	// in: body
	// required: true
	Body PostUpdateRequestBody
}

// PostsFetchRequestQuery includes the ID and creation time of the last seen post,
// as well as a count of the number of posts to retrieve.
type PostsFetchRequestQuery struct {
	// in: query
	LastSeenPostID uuid.UUID `json:"last_seen_post_id" validate:"uuid"`

	// in: query
	LastSeenPostCreatedAt time.Time `json:"last_seen_post_created_at" validate:"uuid,required_with=last_seen_post_id"`

	// in: query
	// required: true
	// min: 1
	// max: 50
	Count int `json:"count" validate:"required,min=1,max=50"`
}

// PostsFetchRequest is a struct that encapsulates a query used to fetch posts.
// swagger:parameters getPosts
type PostsFetchRequest struct {
	PostsFetchRequestQuery
}
