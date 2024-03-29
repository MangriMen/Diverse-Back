package parameters

import (
	"time"

	"github.com/google/uuid"
)

// PostFetchType is type for relations between users.
type PostFetchType string

// Enum for relation type.
const (
	Subscriptions PostFetchType = "subscriptions"
	User          PostFetchType = "user"
	All           PostFetchType = "all"
)

// PostIDParams includes the id of the post.
type PostIDParams struct {
	// in: path
	// required: true
	Post uuid.UUID `params:"post" json:"post" validate:"required"`
}

// PostIDRequest is used to represent a request that requires a post id parameter,
// such as fetching a specific post or deleting a post.
// swagger:parameters getPost updatePost deletePost likePost unlikePost
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

// PostsFetchCountRequestQuery includes the ID and creation time of the last seen post,
// as well as a count of the number of posts to retrieve.
type PostsFetchCountRequestQuery struct {
	// in: query
	Type PostFetchType `query:"type" json:"type" validate:"required"`

	// in: query
	UserID uuid.UUID `query:"user_id" json:"user_id" validate:"uuid,required_with=type"`
}

// PostsFetchCountRequest is a struct that encapsulates a query used to fetch posts count.
// swagger:parameters getPostsCount
type PostsFetchCountRequest struct {
	PostsFetchCountRequestQuery
}

// PostsFetchRequestQuery includes the ID and creation time of the last seen post,
// as well as a count of the number of posts to retrieve.
type PostsFetchRequestQuery struct {
	// in: query
	LastSeenPostID uuid.UUID `query:"last_seen_post_id" json:"last_seen_post_id" validate:"uuid"`

	//nolint:lll
	// in: query
	LastSeenPostCreatedAt time.Time `query:"last_seen_post_created_at" json:"last_seen_post_created_at" validate:"uuid,required_with=last_seen_post_id"`

	// in: query
	Type PostFetchType `query:"type" json:"type" validate:"required"`

	// in: query
	UserID uuid.UUID `query:"user_id" json:"user_id" validate:"uuid,required_with=type"`

	// in: query
	// required: true
	// min: 1
	// max: 50
	Count int `query:"count" json:"count" validate:"required,min=1,max=50"`
}

// PostsFetchRequest is a struct that encapsulates a query used to fetch posts.
// swagger:parameters getPosts
type PostsFetchRequest struct {
	PostsFetchRequestQuery
}
