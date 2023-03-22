package responses

import "github.com/MangriMen/Diverse-Back/internal/models"

type GetPostsResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`

	// required: true
	Posts []models.Post `json:"posts"`
}

// swagger:response
type GetPostsResponse struct {
	// in: body
	Body GetPostsResponseBody
}

type GetPostResponseBody struct {
	BaseResponseBody

	// required: true
	Post models.Post `json:"post"`
}

// swagger:response
type GetPostResponse struct {
	// in: body
	Body GetPostResponseBody
}

// swagger:response
type CreateUpdatePostResponse string

// swagger:response
type DeletePostResponse struct {
}

// swagger:response
type CreateUpdateCommentResponse string

// swagger:response
type DeleteCommentResponse struct {
}
