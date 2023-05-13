package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserQueries is struct for interacting with a database for user-related queries.
type UserQueries struct {
	*sqlx.DB
}

// GetUsers is used to fetch users.
func (q *UserQueries) GetUsers() ([]models.DBUser, error) {
	users := []models.DBUser{}

	query := `SELECT *
		FROM users_view`

	err := q.Select(&users, query)
	if err != nil {
		return users, err
	}

	return users, nil
}

// GetUser retrieves a single user from the database based on the given id parameter.
func (q *UserQueries) GetUser(id uuid.UUID) (models.DBUser, error) {
	user := models.DBUser{}

	query := `SELECT *
		FROM users_view
		WHERE id = $1`

	err := q.Get(&user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByEmail retrieves a single user from the database based on the given email parameter.
func (q *UserQueries) GetUserByEmail(email string) (models.DBUser, error) {
	user := models.DBUser{}

	query := `SELECT *
		FROM users_view
		WHERE email = $1`

	err := q.Get(&user, query, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByUsername retrieves a single user from the database based on the given username parameter.
func (q *UserQueries) GetUserByUsername(username string) (models.DBUser, error) {
	user := models.DBUser{}

	query := `SELECT *
		FROM users_view
		WHERE username = $1`

	err := q.Get(&user, query, username)
	if err != nil {
		return user, err
	}

	return user, nil
}

// CreateUser creates a new user at the database based on the given user object.
func (q *UserQueries) CreateUser(b *models.DBUser) error {
	query := `INSERT INTO users
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(query, b.ID, b.Email, b.Password, b.Username, b.Name, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates user content based on the given ID.
func (q *UserQueries) UpdateUser(b *models.DBUser) error {
	query := `UPDATE users
		SET
			email = $2,
			password = $3,
			username = $4,
			name = $5,
			updated_at = $6,
			avatar_url = $7
		WHERE id = $1`

	_, err := q.Exec(
		query,
		b.ID,
		b.Email,
		b.Password,
		b.Username,
		b.Name,
		b.UpdatedAt,
		b.AvatarURL,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes user based on the given ID.
func (q *UserQueries) DeleteUser(id uuid.UUID) error {
	query := `UPDATE users
		SET 
			deleted_at = now()
		WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
