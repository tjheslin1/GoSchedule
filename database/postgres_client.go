package database

import (
	"database/sql"
	"log"
)

// PostgresDBClient is the Postgresql implementation of DBClient.
type PostgresDBClient struct {
	Logger *log.Logger
	db     *sql.DB
}

// Connection opens a database connection if not already connected, rerturning the connection.
func (pgdbClient *PostgresDBClient) Connection() *sql.DB {
	if pgdbClient.db != nil {
		return pgdbClient.db
	}

	db, err := sql.Open("postgres", "user=go_scheduler dbname=go_scheduler_db sslmode=disable")
	check(err, pgdbClient.Logger)

	pgdbClient.db = db
	return db
}

// SubmitJob inserts the specified job into the database.
func (pgdbClient *PostgresDBClient) SubmitJob(jobEntry TableEntry) error {
	// pgdbClient.Connection().Exec(insertStatement)
	return nil
}
