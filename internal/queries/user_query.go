package queries

import (
	"github.com/MangriMen/Value-Back/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) GetUsers() ([]models.User, error) {
	users := []models.User{}

	query := `SELECT * FROM users`

	err := q.Select(&users, query)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (q *UserQueries) GetUser(id uuid.UUID) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE id = $1`

	err := q.Get(&user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE email = $1`

	err := q.Get(&user, query, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) GetUserByUsername(username string) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE username = $1`

	err := q.Get(&user, query, username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) CreateUser(b *models.User) error {
	query := `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(query, b.Id, b.Email, b.Password, b.Username, b.Name, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) UpdateUser(id uuid.UUID, b *models.User) error {
	query := `UPDATE users
		SET
			email = $2,
			password = $3,
			username = $4,
			name = $5,
			created_at = $6,
			updated_at = $7
		WHERE id = $1`

	_, err := q.Exec(query, id, b.Email, b.Password, b.Username, b.Name, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
