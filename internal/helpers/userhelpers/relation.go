// Package userhelpers provides functionality to work with users.
package userhelpers

import (
	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/models"
)

// PrepareRelationToSend prepares a relation object for sending by fetching additional data from the database
// such as the relation user associated with the relation.
func PrepareRelationToSend(relation models.DBRelation, db *database.Queries) *models.Relation {
	preparedRelation := relation.ToRelation()

	user, err := db.GetUser(relation.RelationUserID)
	if err != nil {
		return nil
	}
	preparedRelation.RelationUser = user.ToUser()

	return &preparedRelation
}

// PrepareRelationStatusToSend transforms the list of relations for the user
// into a structure reflecting their status.
func PrepareRelationStatusToSend(relationStatus []models.DBRelation) map[models.RelationType]bool {
	preparedStatus := make(map[models.RelationType]bool)

	for _, relation := range relationStatus {
		preparedStatus[relation.Type] = true
	}

	return preparedStatus
}
