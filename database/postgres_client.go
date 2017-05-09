package database

import (
	"database/sql"
	"log"
)

// PostgresDBClient is the Postgresql implementation of DBClient.
type PostgresDBClient struct {
	Logger *log.Logger
}

// Connect opens a database connection, rerturning the connection.
func (pgdbClient PostgresDBClient) Connect() *sql.DB {
	db, err := sql.Open("postgres", "user=go_scheduler dbname=go_scheduler_db sslmode=disable")
	check(err, pgdbClient.Logger)

	return db
}

// SubmitJob inserts the specified job into the database.
func (pgdbClient PostgresDBClient) SubmitJob(db *sql.DB, jobEntry TableEntry) error {
	return nil
}
