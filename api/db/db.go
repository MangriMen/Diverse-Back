package db

import "database/sql"

func OpenDBConnection() (*sql.DB, error) {
	db, err := PostgreSQLConnection()

	if err != nil {
		return nil, err
	}

	return db, nil
}
