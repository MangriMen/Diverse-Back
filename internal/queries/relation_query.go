package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// RelationQueries is struct for interacting with a database for relation-related queries.
type RelationQueries struct {
	*sqlx.DB
}

// GetRelationsCount retrieves a count of given relation type between users.
func (q *RelationQueries) GetRelationsCount(
	userID uuid.UUID,
	relationGetRequestQuery *parameters.RelationGetCountRequestQuery,
) (int, error) {
	relationsCount := 0

	query := `SELECT Count(*)
		FROM user_relations
		WHERE user_id = $1
		AND type = $2`

	err := q.Get(&relationsCount, query, userID, relationGetRequestQuery.Type)
	if err != nil {
		return relationsCount, err
	}

	return relationsCount, nil
}

// GetRelations retrieves a list of given relation type between users.
func (q *RelationQueries) GetRelations(
	userID uuid.UUID,
	relationGetRequestQuery *parameters.RelationGetRequestQuery,
) ([]models.DBRelation, error) {
	relations := []models.DBRelation{}

	query := `SELECT *
		FROM user_relations
		WHERE created_at < $1
		AND id <> $2
		AND user_id = $4
		AND type = $5
		ORDER BY created_at DESC
		FETCH FIRST $3 ROWS ONLY`

	err := q.Select(
		&relations,
		query,
		relationGetRequestQuery.LastSeenRelationCreatedAt,
		relationGetRequestQuery.LastSeenRelationID,
		relationGetRequestQuery.Count,
		userID,
		relationGetRequestQuery.Type,
	)
	if err != nil {
		return relations, err
	}

	return relations, nil
}

// GetRelationStatus is used to fetch relation status with given user.
func (q *RelationQueries) GetRelationStatus(
	relationGetStatusParams *parameters.RelationGetStatusParams,
) ([]models.DBRelation, error) {
	relations := []models.DBRelation{}

	query := `SELECT *
		FROM user_relations
		WHERE user_id = $1
		AND relation_user_id = $2`

	err := q.Select(
		&relations,
		query,
		relationGetStatusParams.User,
		relationGetStatusParams.RelationUser,
	)
	if err != nil {
		return relations, err
	}

	return relations, nil
}

// AddRelation is used to add new relation with given parameters.
func (q *RelationQueries) AddRelation(r *models.DBRelation) error {
	query := `INSERT INTO user_relations VALUES ($1, $2, $3, $4, $5)`

	_, err := q.Exec(query, r.ID, r.UserID, r.RelationUserID, r.Type, r.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRelation is used to delete relation with given user and type.
func (q *RelationQueries) DeleteRelation(
	relationGetStatusParams *parameters.RelationGetStatusParams,
	relationAddDeleteRequestBody *parameters.RelationAddDeleteRequestBody,
) error {
	query := `DELETE
				FROM user_relations
				WHERE user_id = $1
				AND relation_user_id = $2
				AND type = $3`

	_, err := q.Exec(
		query,
		relationGetStatusParams.User,
		relationGetStatusParams.RelationUser,
		relationAddDeleteRequestBody.Type,
	)
	if err != nil {
		return err
	}

	return nil
}
