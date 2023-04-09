// Package userhelpers provides functionality to convert relation
// from DB to response variant.
package userhelpers

import (
	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/models"
)

// PrepareRelationToSend prepares a relation object for sending by fetching additional data from the database
// such as the relation user associated with the relation.
// The function then returns the prepared relation object.
func PrepareRelationToSend(relation models.DBRelation, db *database.Queries) *models.Relation {
	preparedRelation := relation.ToRelation()

	user, err := db.GetUser(relation.RelationUserID)
	if err != nil {
		return nil
	}

	preparedRelation.RelationUser = user.ToUser()

	return &preparedRelation
}

// PrepareRelationStatusToSend converts a list of relation for user into
// a structure that reflects their presence.
func PrepareRelationStatusToSend(relationStatus []models.DBRelation) map[models.RelationType]bool {
	preparedStatus := make(map[models.RelationType]bool)

	for _, relation := range relationStatus {
		preparedStatus[relation.Type] = true
	}

	return preparedStatus
}
