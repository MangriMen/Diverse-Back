package database

import (
	"github.com/MangriMen/Diverse-Back/internal/queries"
)

type Queries struct {
	*queries.UserQueries
	*queries.PostQueries
	*queries.CommentQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:    &queries.UserQueries{DB: db},
		PostQueries:    &queries.PostQueries{DB: db},
		CommentQueries: &queries.CommentQueries{DB: db},
	}, nil
}
