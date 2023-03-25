package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostQueries is struct for interacting with a database for post-related queries.
type PostQueries struct {
	*sqlx.DB
}

// GetPosts is used to fetch posts.
// Returns a slice of posts.
func (q *PostQueries) GetPosts(
	postsFetchRequestQuery *parameters.PostsFetchRequestQuery,
) ([]models.DBPost, error) {
	posts := []models.DBPost{}

	query := `SELECT *
		FROM posts
		WHERE created_at < $1
		ORDER BY created_at DESC
		FETCH FIRST $2 ROWS ONLY`

	err := q.Select(
		&posts,
		query,
		postsFetchRequestQuery.LastSeenPostCreatedAt,
		postsFetchRequestQuery.Count,
	)
	if err != nil {
		return posts, err
	}

	return posts, nil
}

// GetPost retrieves a single post from the database based on the given id parameter.
func (q *PostQueries) GetPost(id uuid.UUID) (models.DBPost, error) {
	post := models.DBPost{}

	query := `SELECT * FROM posts WHERE id = $1`

	err := q.Get(&post, query, id)
	if err != nil {
		return post, err
	}

	return post, nil
}

// CreatePost creates a new post at the database based on the given post object.
func (q *PostQueries) CreatePost(b *models.DBPost) error {
	query := `INSERT INTO posts VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.Exec(query, b.ID, b.UserID, b.Content, b.Description, b.Likes, b.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePost updates post content based on the given ID.
func (q *UserQueries) UpdatePost(b *models.DBPost) error {
	query := `UPDATE posts
		SET
			description = $2
		WHERE id = $1`

	_, err := q.Exec(query, b.ID, b.Description)
	if err != nil {
		return err
	}

	return nil
}

// DeletePost deletes post based on the given ID.
func (q *PostQueries) DeletePost(id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
