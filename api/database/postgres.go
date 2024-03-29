package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MangriMen/Diverse-Back/internal/helpers"
	_ "github.com/jackc/pgx/v4/stdlib" // compatibility layer for sqlx
	"github.com/jmoiron/sqlx"
)

// PostgreSQLConnection open connection to postgres database
// with parameters from environment.
func PostgreSQLConnection() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	dataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_TYPE"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	if err = db.Ping(); err != nil {
		defer helpers.CloseQuietly(db)
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
