package queries

import (
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) GetUsers() ([]models.DBUser, error) {
	users := []models.DBUser{}

	query := `SELECT * FROM users`

	err := q.Select(&users, query)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (q *UserQueries) GetUser(id uuid.UUID) (models.DBUser, error) {
	user := models.DBUser{}

	query := `SELECT * FROM users WHERE id = $1`

	err := q.Get(&user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) GetUserByEmail(email string) (models.DBUser, error) {
	user := models.DBUser{}

	query := `SELECT * FROM users WHERE email = $1`

	err := q.Get(&user, query, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) GetUserByUsername(username string) (models.DBUser, error) {
	user := models.DBUser{}

	query := `SELECT * FROM users WHERE username = $1`

	err := q.Get(&user, query, username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) CreateUser(b *models.DBUser) error {
	query := `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(query, b.Id, b.Email, b.Password, b.Username, b.Name, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) UpdateUser(b *models.DBUser) error {
	query := `UPDATE users
		SET
			email = $2,
			password = $3,
			username = $4,
			name = $5,
			updated_at = $6
		WHERE id = $1`

	_, err := q.Exec(query, b.Id, b.Email, b.Password, b.Username, b.Name, b.UpdatedAt)
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
