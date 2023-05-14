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

// GetCommentsCount is used to fetch comments count.
func (q *PostQueries) GetCommentsCount(postID uuid.UUID) (int, error) {
	commentsCount := 0

	query := `SELECT Count(*)
		FROM comments_view
		WHERE post_id = $1`

	err := q.Get(
		&commentsCount,
		query,
		postID,
	)
	if err != nil {
		return commentsCount, err
	}

	return commentsCount, nil
}

// GetComments is used to fetch comments related to a post based
// on a provided post ID.
func (q *PostQueries) GetComments(
	postID uuid.UUID,
	commentsFetchRequestQuery *parameters.CommentsFetchRequestQuery,
) ([]models.DBComment, error) {
	comments := []models.DBComment{}

	query := `SELECT *
		FROM comments_view
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

	query := `SELECT *
		FROM comments_view
		WHERE id = $1`

	err := q.Get(&comment, query, id)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// AddComment add a single comment to the database based on the given comment object.
func (q *PostQueries) AddComment(b *models.DBComment) error {
	query := `INSERT INTO comments
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(query, b.ID, b.PostID, b.UserID, b.Content, b.CreatedAt, b.UpdatedAt, b.Likes)
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

// LikeComment sets like the comment by ID.
func (q *PostQueries) LikeComment(l *models.DBCommentLike) error {
	query := `INSERT INTO comment_likes
		VALUES ($1, $2, $3)`

	_, err := q.Exec(query, l.ID, l.CommentID, l.UserID)
	if err != nil {
		return err
	}

	return nil
}

// UnlikeComment sets like the comment by ID.
func (q *PostQueries) UnlikeComment(l *models.DBCommentLike) error {
	query := `DELETE
		FROM comment_likes
		WHERE comment_id = $1 AND user_id = $2`

	_, err := q.Exec(query, l.CommentID, l.UserID)
	if err != nil {
		return err
	}

	return nil
}

// GetCommentIsLiked gets status of like for given comment and user.
func (q *PostQueries) GetCommentIsLiked(commentID uuid.UUID, userID uuid.UUID) (bool, error) {
	likesCount := 0

	query := `SELECT Count(*)
		FROM comment_likes
		WHERE comment_id = $1 AND user_id = $2`

	err := q.Get(&likesCount, query, commentID, userID)
	if err != nil {
		return false, err
	}

	return likesCount > 0, nil
}

// DeleteComment deletes comment from database based on the given comment ID.
func (q *PostQueries) DeleteComment(id uuid.UUID) error {
	query := `UPDATE comments
		SET
			deleted_at = now()
		WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
