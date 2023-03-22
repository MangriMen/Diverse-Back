package responses

import "github.com/MangriMen/Diverse-Back/internal/models"

type GetCommentsResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`

	// required: true
	Comments []models.Comment `json:"comments"`
}

// swagger:response
type GetCommentsResponse struct {
	// in: body
	Body GetCommentsResponseBody
}

type GetCommentResponseBody struct {
	BaseResponseBody

	// required: true
	Comment models.Comment `json:"comment"`
}

// swagger:response
type GetCommentResponse struct {
	// in: body
	Body GetCommentResponseBody
}
