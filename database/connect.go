package database

import (
	"database/sql"
	"log"
)

// Connect opens a database connection, rerturning the connection.
func Connect(logger *log.Logger) *sql.DB {
	db, err := sql.Open("postgres", "user=go_scheduler dbname=go_scheduler_db sslmode=disable")
	check(err, logger)

	return db
}

func check(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
