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
		_, err := db.Exec(`CREATE TABLE jobs(
			JOB_ID INT PRIMARY KEY NOT NULL,
			NAME TEXT NOT NULL);`)
		check(err, logger)

		logger.Println("'jobs' table created")
	}
}

func check(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
