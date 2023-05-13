package responses

import "github.com/MangriMen/Diverse-Back/internal/models"

// GetPostsCountResponseBody includes the posts count.
type GetPostsCountResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`
}

// GetPostsCountResponse represent the response retrived on get posts request.
// swagger:response
type GetPostsCountResponse struct {
	// in: body
	Body GetPostsCountResponseBody
}

// GetPostsResponseBody includes the slice of posts.
type GetPostsResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`

	// required: true
	Data []models.Post `json:"data"`
}

// GetPostsResponse represent the response retrived on get posts request.
// swagger:response
type GetPostsResponse struct {
	// in: body
	Body GetPostsResponseBody
}

// GetPostResponseBody includes the signle post for a given ID.
type GetPostResponseBody struct {
	BaseResponseBody

	// required: true
	Data models.Post `json:"data"`
}

// GetPostResponse represent the response retrived on get post request.
// swagger:response
type GetPostResponse struct {
	// in: body
	Body GetPostResponseBody
}

// CreateUpdatePostResponse represents response for successfully create or update post request.
// swagger:response
type CreateUpdatePostResponse string

// DeletePostResponse represents response for successfully delete post.
// swagger:response
type DeletePostResponse struct {
}

// CreateUpdateCommentResponse represents response for successfully create or update comment request.
// swagger:response
type CreateUpdateCommentResponse string

// DeleteCommentResponse represents response for successfully delete comment request.
// swagger:response
type DeleteCommentResponse struct {
}
