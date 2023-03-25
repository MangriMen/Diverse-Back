package helpers

import "github.com/jmoiron/sqlx"

// CloseDBQuietly close db object quietly without returning error.
func CloseDBQuietly(db *sqlx.DB) {
	if err := db.Close(); err != nil {
		_ = db.Close()
	}
}
