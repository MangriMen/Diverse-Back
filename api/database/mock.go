package database

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

// InitUsersMock inits mock data for users.
func InitUsersMock(mock sqlmock.Sqlmock) {
	columns := []string{
		"id",
		"email",
		"password",
		"username",
		"name",
		"created_at",
		"updated_at",
		"avatar_url",
	}

	rows1 := sqlmock.NewRows(columns).AddRow(
		"17afe1a4-aaed-4263-a29b-781389509cb6",
		"test@test.test",
		"$2a$12$sk5KosZRDr.I3egLJJd7rujP/6lG69mC.OZGKLIqgCWpYUVYDYQ.W",
		"test_user",
		"",
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
		"/data/image/avatar",
	)

	rows2 := sqlmock.NewRows(columns).AddRow(
		"17afe1a4-aaed-4263-a29b-781389509cb6",
		"test@test.test",
		"$2a$12$sk5KosZRDr.I3egLJJd7rujP/6lG69mC.OZGKLIqgCWpYUVYDYQ.W",
		"test_user",
		"",
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
		"/data/image/avatar",
	)

	mock.ExpectQuery(`SELECT * FROM users_view WHERE email = $1`).
		WithArgs("test@test.test").
		WillReturnRows(rows1)

	mock.ExpectQuery(`SELECT * FROM users_view`).
		WillReturnRows(rows2)
}

// InitPostsMock inits mock data for posts.
func InitPostsMock(mock sqlmock.Sqlmock) {
	columns := []string{
		"id",
		"user_id",
		"content",
		"description",
		"likes",
		"created_at",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		"023753e6-6539-4c95-93e6-1d469d80f28d",
		"17afe1a4-aaed-4263-a29b-781389509cb6",
		"data/image/image0",
		"test description",
		0,
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
	)

	count := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)

	mock.ExpectQuery(`SELECT Count(*) FROM posts_view WHERE 1 = 1`).
		WillReturnRows(count)

	mock.ExpectQuery(`SELECT *
		FROM posts_view
		WHERE created_at < $1
		AND id <> $2
		ORDER BY created_at DESC
		FETCH FIRST $3 ROWS ONLY`).
		WithArgs(
			time.Date(2024, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
			"",
			10,
		).
		WillReturnRows(rows)
}

// InitCommentsMock inits mock data for comments.
func InitCommentsMock(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"id",
		"post_id",
		"user_id",
		"content",
		"created_at",
		"updated_at",
		"likes",
	}).AddRow(
		"3c0c325a-d2f0-4a9d-ac8e-6b791509f162",
		"023753e6-6539-4c95-93e6-1d469d80f28d",
		"17afe1a4-aaed-4263-a29b-781389509cb6",
		"test comment",
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
		time.Date(2023, time.Month(1), 2, 3, 4, 5, 6, time.UTC),
		0,
	)

	count := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)

	mock.ExpectQuery(`SELECT Count(*) FROM comments_view WHERE post_id = $1`).
		WithArgs("023753e6-6539-4c95-93e6-1d469d80f28d").
		WillReturnRows(count)

	mock.ExpectQuery(`SELECT *
			FROM comments_view
			WHERE post_id = $1
			AND created_at < $2
			ORDER BY created_at DESC
			FETCH FIRST $3 ROWS ONLY`).
		WithArgs("023753e6-6539-4c95-93e6-1d469d80f28d", time.Date(2024, time.Month(1), 2, 3, 4, 5, 6, time.UTC), 10).
		WillReturnRows(rows)
}

// InitMock inits all mock queries.
func InitMock(mock sqlmock.Sqlmock) {
	InitUsersMock(mock)
	InitPostsMock(mock)
	InitCommentsMock(mock)
}

// MockSQLConnection create mock database.
func MockSQLConnection() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s_%d", configs.MockDSN, rand.Intn(10000))

	db, mock, err := sqlmock.NewWithDSN(dsn, sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	if os.Getenv("ENABLE_MOCK_LOGS") != "" {
		db = sqldblogger.OpenDriver(
			dsn,
			db.Driver(),
			zerologadapter.New(zerolog.New(helpers.OpenTempFile(configs.MockLogFilename))),
		)
	}

	if err = db.Ping(); err != nil {
		defer helpers.CloseQuietly(db)
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	InitMock(mock)

	sqlxDB := sqlx.NewDb(db, "mock")

	return sqlxDB, nil
}
