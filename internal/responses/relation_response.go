package responses

import "github.com/MangriMen/Diverse-Back/internal/models"

// GetRelationResponseBody includes the slice of relations.
type GetRelationResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`

	// required: true
	Relations []models.Relation `json:"relations"`
}

// GetRelationResponse represent the response retrived on get relation request.
// swagger:response
type GetRelationResponse struct {
	// in: body
	Body GetRelationResponseBody
}

// AddRelationResponse represents response for successfully create relation.
// swagger:response
type AddRelationResponse string

// GetRelationStatusResponseBody includes the status of relations.
type GetRelationStatusResponseBody struct {
	Follower  bool `json:"follower"`
	Following bool `json:"following"`
	Blocked   bool `json:"blocked"`
}

// GetRelationStatusResponse represent the response retrived on
// get relation by relation user request.
// swagger:response
type GetRelationStatusResponse struct {
	GetRelationStatusResponseBody
}

// DeleteRelationResponse represents response for successfully delete relation request.
// swagger:response
type DeleteRelationResponse struct {
}
