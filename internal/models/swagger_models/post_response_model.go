package swagger_models

import "github.com/MangriMen/Diverse-Back/internal/models"

// swagger:response
type GetPostsResponse struct {
	// in: body
	Body struct {
		BaseResponseBody
		// required: true
		Posts []models.Post `json:"users"`
		// required: true
		Count int `json:"count"`
	}
}

// swagger:response
type GetPostResponse struct {
	// in: body
	Body struct {
		BaseResponseBody
		// required: true
		Post models.Post `json:"user"`
	}
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
