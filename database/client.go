package database

import (
	"database/sql"
	"log"
)

// TableEntry represents data to be inserted into the database,
// the Data field is a map of column name -> value for a new row.
type TableEntry struct {
	Name string
	Data map[string]interface{}
}

// DBClient defines the interactions with the database.
type DBClient interface {
	Connect() *sql.DB
	SubmitJob(*sql.DB, TableEntry) bool
}

// PostgresDBClient is the Postgresql implementation of DBClient.
type PostgresDBClient struct {
	Logger *log.Logger
}

// Connect opens a database connection, rerturning the connection.
func (pgdbClient *PostgresDBClient) Connect() *sql.DB {
	db, err := sql.Open("postgres", "user=go_scheduler dbname=go_scheduler_db sslmode=disable")
	check(err, pgdbClient.Logger)

	return db
}

// SubmitJob inserts the specified job into the database.
func (pgdbClient *PostgresDBClient) SubmitJob(db *sql.DB, jobEntry TableEntry) error {
	return nil
}

func check(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
