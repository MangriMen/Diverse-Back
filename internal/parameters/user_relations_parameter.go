package parameters

import (
	"time"

	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/google/uuid"
)

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
// swagger:parameters getRelation
type RelationGetRequest struct {
	RelationGetRequestQuery
}
