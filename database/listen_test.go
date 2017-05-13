package database

import (
	"bytes"
	"database/sql"
	"log"
	"testing"
)

func TestWatch(t *testing.T) {
	// taskListener := taskListener{newDummyLogger{}, dummyDBClient{}}
	taskListener := taskListener{newDummyLogger(), &PostgresDBClient{}}

	taskListener.watch()
}

func newDummyLogger() *log.Logger {
	return log.New(new(bytes.Buffer), "", 0)
}

type dummyDBClient struct {
	capturedTableEntry TableEntry
}

func (dummyClient *dummyDBClient) Connection() *sql.DB {
	return nil
}

func (dummyClient *dummyDBClient) SubmitEntry(jobEntry TableEntry) error {
	dummyClient.capturedTableEntry = jobEntry
	return nil
}
