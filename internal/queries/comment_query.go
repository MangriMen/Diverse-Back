package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentQueries struct {
	*sqlx.DB
}

func (q *PostQueries) GetComments(postId uuid.UUID) ([]models.DBComment, error) {
	comments := []models.DBComment{}

	query := `SELECT * FROM comments WHERE post_id = $1`

	err := q.Select(&comments, query, postId)
	if err != nil {
		return comments, err
	}

	return comments, nil
}

func (q *PostQueries) GetComment(id uuid.UUID) (models.DBComment, error) {
	comment := models.DBComment{}

	query := `SELECT * FROM comments WHERE id = $1`

	err := q.Get(&comment, query, id)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (q *PostQueries) AddComment(b *models.DBComment) error {
	query := `INSERT INTO comments VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.Exec(query, b.Id, b.PostId, b.UserId, b.Content, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) UpdateComment(id uuid.UUID, b *models.DBComment) error {
	query := `UPDATE comments
		SET
			content = $2
			updated_at = $3
		WHERE id = $1`

	_, err := q.Exec(query, id, b.Content, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (q *PostQueries) DeleteComment(id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
