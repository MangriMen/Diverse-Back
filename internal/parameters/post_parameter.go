package parameters

import (
	"time"

	"github.com/google/uuid"
)

type PostIdParams struct {
	// in: path
	// required: true
	Post uuid.UUID `params:"post" json:"post" validate:"required"`
}

// swagger:parameters getPost updatePost deletePost addComment updateComment deleteComment
type PostIdRequest struct {
	PostIdParams
}

type PostCreateRequestBody struct {
	// required: true
	Content string `json:"content" validate:"required"`

	// required: true
	// max length: 2048
	Description string `json:"description" validate:"lte=2048"`

	// required: true
	UserId uuid.UUID `json:"user_id" validate:"required,uuid"`
}

// swagger:parameters createPost
type PostCreateRequest struct {
	// in: body
	// required: true
	Body PostCreateRequestBody
}

type PostUpdateRequestBody struct {
	// required: true
	// max length: 2048
	Description string `json:"description" validate:"lte=2048"`
}

type PostUpdateRequest struct {
	PostIdParams

	// in: body
	// required: true
	Body PostUpdateRequestBody
}

type PostsFetchRequestQuery struct {
	// in: query
	LastSeenPostId uuid.UUID `json:"last_seen_post_id" validate:"uuid"`

	// in: query
	LastSeenPostCreatedAt time.Time `json:"last_seen_post_created_at" validate:"uuid,required_with=last_seen_post_id"`

	// in: query
	// required: true
	// min: 1
	// max: 50
	Count int `json:"count" validate:"required,min=1,max=50"`
}

// swagger:parameters getPosts
type PostsFetchRequest struct {
	PostsFetchRequestQuery
}
