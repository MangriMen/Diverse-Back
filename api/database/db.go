// Package database provides database connection and convenient queries
package database

import (
	"os"

	"github.com/MangriMen/Diverse-Back/internal/queries"
	"github.com/jmoiron/sqlx"
)

// Queries is struct for storing various queries.
type Queries struct {
	*queries.UserQueries
	*queries.RelationQueries
	*queries.PostQueries
	*queries.CommentQueries
}

// OpenDBConnection open db connection and combine all queries.
func OpenDBConnection() (*Queries, error) {
	var db *sqlx.DB
	var err error

	if os.Getenv("ENABLE_TESTING") != "" {
		db, err = MockSQLConnection()
	} else {
		db, err = PostgreSQLConnection()
	}

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries:     &queries.UserQueries{DB: db},
		RelationQueries: &queries.RelationQueries{DB: db},
		PostQueries:     &queries.PostQueries{DB: db},
		CommentQueries:  &queries.CommentQueries{DB: db},
	}, nil
}
