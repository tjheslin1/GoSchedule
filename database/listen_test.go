package database

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"testing"
)

func TestReportProblemCallback(t *testing.T) {
	dummyLogger := newDummyLogger()

	taskListener := JobListener{dummyLogger.logger}

	taskListener.reportProblemCallback(0, errors.New("Test error"))

	var expectedLogOuput = "Error occured for pq.ListenerEventType: '0'.\nTest error\n"
	if dummyLogger.logOutput() != expectedLogOuput {
		t.Errorf("Expected log output of:'%s'\nbut got:\n'%v'\n", expectedLogOuput, dummyLogger.logOutput())
	}
}

type dummyLogger struct {
	logger    *log.Logger
	logBuffer bytes.Buffer
}

func newDummyLogger() *dummyLogger {
	dummyLogger := dummyLogger{}

	dummyLogger.logger = log.New(nil, "", 0)
	dummyLogger.logger.SetOutput(&dummyLogger.logBuffer)

	return &dummyLogger
}

func (dl *dummyLogger) logOutput() string {
	return dl.logBuffer.String()
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
