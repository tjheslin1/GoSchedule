package database

import (
	"database/sql"
	"log"
)

var connectionInfo = "user=go_scheduler dbname=go_scheduler_db sslmode=disable"

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

	db, err := sql.Open("postgres", connectionInfo)
	check(err, pgdbClient.Logger)

	pgdbClient.db = db
	return db
}

// SubmitEntry inserts the specified job into the database.
func (pgdbClient *PostgresDBClient) SubmitEntry(jobEntry TableEntry) error {
	statement := insertStatement(jobEntry)
	pgdbClient.Logger.Printf("Executing query: '%s'", statement)
	_, err := pgdbClient.Connection().Exec(statement)
	return err
}
