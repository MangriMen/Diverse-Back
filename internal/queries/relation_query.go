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

// GetRelations retrieves a list of given relation type between users.
func (q *RelationQueries) GetRelations(
	userID uuid.UUID,
	relationGetRequestQuery *parameters.RelationGetRequestQuery,
) ([]models.DBRelation, error) {
	relations := []models.DBRelation{}

	query := `SELECT *
		FROM user_relations
		WHERE created_at < $1
		AND user_id = $3
		AND type = $4
		ORDER BY created_at DESC
		FETCH FIRST $2 ROWS ONLY`

	err := q.Select(
		&relations,
		query,
		relationGetRequestQuery.LastSeenRelationCreatedAt,
		relationGetRequestQuery.Count,
		userID,
		relationGetRequestQuery.Type,
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
