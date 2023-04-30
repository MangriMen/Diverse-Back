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

	const baseQuery = `SELECT Count(*) FROM user_relations`

	const followingCondition = `WHERE user_id = $1 AND type = $2`
	const followerCondition = `WHERE relation_user_id = $1 AND type = $2`
	const blockedCondition = followingCondition

	relationType := relationGetRequestQuery.Type

	query := baseQuery + ` `

	switch relationGetRequestQuery.Type {
	case models.Follower:
		query += followerCondition
		relationType = models.Following
	case models.Following:
		query += followingCondition
	case models.Blocked:
		query += blockedCondition
	}

	err := q.Get(&relationsCount, query, userID, relationType)
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

	const baseQuery = `SELECT *
		FROM user_relations
		WHERE created_at < $1
		AND id <> $2`

	const cutFilter = `ORDER BY created_at DESC
		FETCH FIRST $3 ROWS ONLY`

	const followingCondition = `AND user_id = $4 AND type = $5`
	const followerCondition = `AND relation_user_id = $4 AND type = $5`
	const blockedCondition = followingCondition

	relationType := relationGetRequestQuery.Type

	query := baseQuery + ` `

	switch relationGetRequestQuery.Type {
	case models.Follower:
		query += followerCondition
		relationType = models.Following
	case models.Following:
		query += followingCondition
	case models.Blocked:
		query += blockedCondition
	}

	query += ` ` + cutFilter

	err := q.Select(
		&relations,
		query,
		relationGetRequestQuery.LastSeenRelationCreatedAt,
		relationGetRequestQuery.LastSeenRelationID,
		relationGetRequestQuery.Count,
		userID,
		relationType,
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
		AND relation_user_id = $2
		OR user_id = $2
		AND relation_user_id = $1`

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
	relationAddDeleteRequestBody *parameters.RelationAddDeleteRequestQuery,
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
