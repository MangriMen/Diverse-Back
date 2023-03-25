// Package queries provides sql queries for database
package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// CommentQueries is struct for interacting with a database for comment-related queries.
type CommentQueries struct {
	*sqlx.DB
}

// GetComments is used to fetch comments related to a post based
// on a provided post ID.
// Returns a slice of comments and a error.
func (q *PostQueries) GetComments(
	postID uuid.UUID,
	commentsFetchRequestQuery *parameters.CommentsFetchRequestQuery,
) ([]models.DBComment, error) {
	comments := []models.DBComment{}

	query := `SELECT *
		FROM comments
		WHERE post_id = $1
		AND created_at < $2
		ORDER BY created_at DESC
		FETCH FIRST $3 ROWS ONLY`

	err := q.Select(
		&comments,
		query,
		postID,
		commentsFetchRequestQuery.LastSeenCommentCreatedAt,
		commentsFetchRequestQuery.Count,
	)
	if err != nil {
		return comments, err
	}

	return comments, nil
}

// GetComment retrieves a single comment from the database based on the given id parameter.
func (q *PostQueries) GetComment(id uuid.UUID) (models.DBComment, error) {
	comment := models.DBComment{}

	query := `SELECT * FROM comments WHERE id = $1`

	err := q.Get(&comment, query, id)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// AddComment add a single comment to the database based on the given comment object.
func (q *PostQueries) AddComment(b *models.DBComment) error {
	query := `INSERT INTO comments VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.Exec(query, b.ID, b.PostID, b.UserID, b.Content, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// UpdateComment updates comment content based on the given comment ID.
func (q *UserQueries) UpdateComment(b *models.DBComment) error {
	query := `UPDATE comments
		SET
			content = $2
			updated_at = $3
		WHERE id = $1`

	_, err := q.Exec(query, b.ID, b.Content, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment deletes comment from database based on the given comment ID.
func (q *PostQueries) DeleteComment(id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
