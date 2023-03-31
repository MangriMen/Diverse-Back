package responses

import "github.com/MangriMen/Diverse-Back/internal/models"

// GetRelationResponseBody includes the slice of relations.
type GetRelationResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`

	// required: true
	Relations []models.Relation `json:"relation"`
}

// GetRelationResponse represent the response retrived on get relation request.
// swagger:response
type GetRelationResponse struct {
	// in: body
	Body GetRelationResponseBody
}
