package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostQueries struct {
	*sqlx.DB
}

func (q *PostQueries) GetPosts() ([]models.DBPost, error) {
	posts := []models.DBPost{}

	query := `SELECT * FROM posts`

	err := q.Select(&posts, query)
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func (q *PostQueries) GetPost(id uuid.UUID) (models.DBPost, error) {
	post := models.DBPost{}

	query := `SELECT * FROM posts WHERE id = $1`

	err := q.Get(&post, query, id)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (q *PostQueries) CreatePost(b *models.DBPost) error {
	query := `INSERT INTO posts VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.Exec(query, b.Id, b.UserId, b.Content, b.Description, b.Likes, b.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) UpdatePost(id uuid.UUID, b *models.DBPost) error {
	query := `UPDATE posts
		SET
			description = $2
		WHERE id = $1`

	_, err := q.Exec(query, id, b.Description)
	if err != nil {
		return err
	}

	return nil
}

func (q *PostQueries) DeletePost(id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
