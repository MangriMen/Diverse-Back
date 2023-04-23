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

// GetRelationsResponse represent the response retrived on get relation request.
// swagger:response
type GetRelationsResponse struct {
	// in: body
	Body GetRelationResponseBody
}

// GetRelationCountResponseBody includes the count of relations.
type GetRelationCountResponseBody struct {
	BaseResponseBody

	// required: true
	Count int `json:"count"`
}

// GetRelationCountResponse represent the response retrived on get relation count request.
// swagger:response
type GetRelationCountResponse struct {
	// in: body
	Body GetRelationCountResponseBody
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
