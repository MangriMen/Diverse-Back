package database

import (
	"github.com/MangriMen/Diverse-Back/internal/queries"
)

type Queries struct {
	*queries.UserQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
	}, nil
}
