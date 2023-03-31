package models

import "github.com/google/uuid"

// RelationType is type for relations between users.
type RelationType string

// Enum for relation type.
const (
	Following RelationType = "following"
	Follower  RelationType = "follower"
	Blocked   RelationType = "blocked"
)

// BaseRelation represents a base relation struct in a system.
type BaseRelation struct {
	// The id for this relation
	// required: true
	ID uuid.UUID `db:"id" json:"id" validate:"required,uuid"`

	// Relation type
	// required: true
	Type RelationType `db:"type" json:"type" validate:"required"`
}

// DBRelation represents a relations struct from database.
type DBRelation struct {
	BaseRelation

	// The id of the user
	// required: true
	UserID uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`

	// The id of the user id between which the relationship
	// required: true
	RelationUserID uuid.UUID `db:"relation_user_id" json:"relation_user_id" validate:"required,uuid"`
}

// ToRelation converts the DBRelation to Relation model.
func (r *DBRelation) ToRelation() Relation {
	return Relation{BaseRelation: r.BaseRelation}
}

// Relation represents the relation for this application
// swagger:model
type Relation struct {
	BaseRelation

	RelationUser User `json:"relation_user"`
}
