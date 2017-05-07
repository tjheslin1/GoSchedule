package database

import (
	"database/sql"
	"log"
)

// SetUpSchema creates any necessary tables if they do not already exist.
func SetUpSchema(db *sql.DB, logger *log.Logger) {
	rows, err := db.Query(`SELECT EXISTS(
	    SELECT *
	    FROM information_schema.tables
	    WHERE
	      table_schema = 'public' AND
	      table_name = 'jobs'
	);`)
	check(err, logger)

	var exists bool
	if rows.Next() {
		err := rows.Scan(&exists)
		check(err, logger)

		logger.Printf("Querying if jobs table exists: '%v'", exists)
	}

	if !exists {
		_, err := db.Exec(createJobsTable)
		check(err, logger)

		logger.Println("'jobs' table created")
	}
}

const createJobsTable string = `CREATE TABLE jobs(
    JOB_ID INT PRIMARY KEY NOT NULL,
    NAME TEXT NOT NULL,
    URL TEXT NOT NULL,
    INTERVAL BIGINT NOT NULL,
    CREATED_AT TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc'));`
