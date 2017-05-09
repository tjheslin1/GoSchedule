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
	Connection() *sql.DB
	SubmitJob(TableEntry) error
}

func check(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
