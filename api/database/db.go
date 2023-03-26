// Package database provides database connection and convenient queries
package database

import (
	"github.com/MangriMen/Diverse-Back/internal/queries"
)

// Queries is struct for storing various queries.
type Queries struct {
	*queries.UserQueries
	*queries.PostQueries
	*queries.CommentQueries
}

// OpenDBConnection open db connection, combine all queries and returns Queries object.
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
