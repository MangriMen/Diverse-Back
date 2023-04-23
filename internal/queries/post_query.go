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

// GetPostsCount is used to fetch posts count.
func (q *PostQueries) GetPostsCount(postFromCondition string) (int, error) {
	postsCount := 0

	query := `SELECT Count(*)
		FROM posts
		WHERE 1 = 1` +
		"\n" + postFromCondition + "\n"

	err := q.Get(
		&postsCount,
		query,
	)
	if err != nil {
		return postsCount, err
	}

	return postsCount, nil
}

// GetPosts is used to fetch posts.
func (q *PostQueries) GetPosts(
	postsFetchRequestQuery *parameters.PostsFetchRequestQuery,
	postFromCondition string,
) ([]models.DBPost, error) {
	posts := []models.DBPost{}

	query := `SELECT *
		FROM posts
		WHERE created_at < $1
		AND id <> $2` +
		"\n" + postFromCondition + "\n" +
		`ORDER BY created_at DESC
		FETCH FIRST $3 ROWS ONLY`

	err := q.Select(
		&posts,
		query,
		postsFetchRequestQuery.LastSeenPostCreatedAt,
		postsFetchRequestQuery.LastSeenPostID,
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

// LikePost sets like the post by ID.
func (q *PostQueries) LikePost(l *models.DBPostLike) error {
	query := `INSERT INTO post_likes VALUES ($1, $2, $3)`

	_, err := q.Exec(query, l.ID, l.PostID, l.UserID)
	if err != nil {
		return err
	}

	return nil
}

// UnlikePost sets like the post by ID.
func (q *PostQueries) UnlikePost(l *models.DBPostLike) error {
	query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`

	_, err := q.Exec(query, l.PostID, l.UserID)
	if err != nil {
		return err
	}

	return nil
}

// GetPostIsLiked gets status of like for given post and user.
func (q *PostQueries) GetPostIsLiked(postID uuid.UUID, userID uuid.UUID) (bool, error) {
	likesCount := 0

	query := `SELECT Count(*) FROM post_likes WHERE post_id = $1 AND user_id = $2`

	err := q.Get(&likesCount, query, postID, userID)
	if err != nil {
		return false, err
	}

	return likesCount > 0, nil
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
