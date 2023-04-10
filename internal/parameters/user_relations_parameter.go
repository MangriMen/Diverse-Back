package parameters

import (
	"time"

	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/google/uuid"
)

// RelationUserIDParams includes the id of the relation user.
type RelationUserIDParams struct {
	// in: path
	// required: true
	RelationUser uuid.UUID `params:"relationUser" json:"relation_user"`
}

// RelationUserIDRequest is used to represent a request that requires a relation user id parameter,
// such as fetching a specific relation for user or deleting a relation .
// swagger:parameters
type RelationUserIDRequest struct {
	RelationUserIDParams
}

// RelationGetRequestQuery includes the ID of the user and relation type to fetch.
type RelationGetRequestQuery struct {
	// in: query
	LastSeenRelationID uuid.UUID `query:"last_seen_relation_id" json:"last_seen_relation_id" validate:"uuid"`

	//nolint:lll
	// in: query
	LastSeenRelationCreatedAt time.Time `query:"last_seen_relation_created_at" json:"last_seen_relation_created_at" validate:"uuid,required_with=last_seen_relation_id"`

	// in: query
	// required: true
	Type models.RelationType `query:"type" json:"type" validate:"required"`

	// in: query
	// required: true
	// min: 1
	// max: 50
	Count int `query:"count" json:"count" validate:"required,min=1,max=50"`
}

// RelationGetRequest is a struct that encapsulates a query used to fetch relation.
// swagger:parameters getRelations
type RelationGetRequest struct {
	RelationGetRequestQuery
}

// RelationGetStatusParams includes the ID of the user.
type RelationGetStatusParams struct {
	UserIDParams
	RelationUserIDParams
}

// RelationGetStatusRequest is a struct that encapsulates a query used to
// fetch relation with specific user.
// swagger:parameters getRelationStatus
type RelationGetStatusRequest struct {
	RelationGetStatusParams
}

// RelationAddDeleteRequestBody includes the ID of the user to add.
type RelationAddDeleteRequestBody struct {
	// Relation type
	// required: true
	Type models.RelationType `json:"type" validate:"required"`
}

// RelationAddDeleteRequest is a struct that encapsulates a body used to add relation.
// swagger:parameters getRelations deleteRelations
type RelationAddDeleteRequest struct {
	UserIDParams
	RelationAddDeleteRequestBody
}