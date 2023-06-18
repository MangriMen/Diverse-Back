// Package responses provides base response models for swagger and work with REST API
package responses

import "github.com/MangriMen/Diverse-Back/internal/models"

// GetCommentsCountResponseBody includes the comments count.
type GetCommentsCountResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`
}

// GetCommentsCountResponse represent the response retrived on get comments request.
// swagger:response
type GetCommentsCountResponse struct {
	// in: body
	Body GetCommentsCountResponseBody
}

// GetCommentsResponseBody includes the slice of comments for a post
// and its count.
type GetCommentsResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`

	// required: true
	Data []models.Comment `json:"data"`
}

// GetCommentsResponse represent the response retrived on get comments request.
// swagger:response
type GetCommentsResponse struct {
	// in: body
	Body GetCommentsResponseBody
}

// GetCommentResponseBody includes the signle comment for a post
// based on given post and comment ID.
type GetCommentResponseBody struct {
	BaseResponseBody

	// required: true
	Data models.Comment `json:"data"`
}

// GetCommentResponse represent the response retrived on get comment request.
// swagger:response
type GetCommentResponse struct {
	// in: body
	Body GetCommentResponseBody
}
