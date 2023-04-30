// Package userhelpers provides functionality to work with users.
package userhelpers

import (
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/google/uuid"
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
func PrepareRelationStatusToSend(userID uuid.UUID, relationStatus []models.DBRelation) map[models.RelationType]bool {
	preparedStatus := make(map[models.RelationType]bool)

	for _, relation := range relationStatus {
		if relation.Type != models.Blocked && userID == relation.RelationUserID {
			preparedStatus[models.Follower] = true
		} else {
			preparedStatus[relation.Type] = true
		}
	}

	return preparedStatus
}

// AddRelation adds relation.
func AddRelation(
	params *parameters.RelationGetStatusParams,
	query *parameters.RelationAddDeleteRequestQuery,
	db *database.Queries) error {
	relation := &models.DBRelation{
		BaseRelation: models.BaseRelation{
			ID:        uuid.New(),
			Type:      query.Type,
			CreatedAt: time.Now(),
		},
		UserID:         params.User,
		RelationUserID: params.RelationUser,
	}

	if err := db.AddRelation(relation); err != nil {
		return err
	}

	return nil
}

// DeleteRealtionWithReverse adds two relation reverse to each other.
func DeleteRealtionWithReverse(
	params *parameters.RelationGetStatusParams,
	query *parameters.RelationAddDeleteRequestQuery,
	db *database.Queries) error {
	if err := db.DeleteRelation(params, query); err != nil {
		return err
	}

	return nil
}
