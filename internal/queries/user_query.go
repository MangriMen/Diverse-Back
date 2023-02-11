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
