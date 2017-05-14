package database

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/tjheslin1/GoSchedule/model"
)

func TestReportProblemCallback(t *testing.T) {
	dummyLogger := newDummyLogger()

	jobListener := JobListener{dummyLogger.logger}

	jobListener.reportProblemCallback(0, errors.New("Test error"))

	var expectedLogOuput = "Error occured for pq.ListenerEventType: '0'.\nTest error\n"
	if dummyLogger.logOutput() != expectedLogOuput {
		t.Errorf("Expected log output of:'%s'\nbut got:\n'%v'\n", expectedLogOuput, dummyLogger.logOutput())
	}
}

func TestUnmarshallJSONJobNotification(t *testing.T) {
	dummyLogger := newDummyLogger()

	notification := []byte(`{
			"table" : "jobs",
			"action" : "INSERT",
			"data" : {
				"job_id":3,
				"name":"testJob",
				"url":"http:localhost:6060/ready",
				"start_time":2,
				"interval":1000,
				"created_at":"2017-05-14T12:25:16.72843+01:00"
			}
		}`)

	jobListener := JobListener{dummyLogger.logger}
	actualSubmitJob := jobListener.unmarshallJSONJobNotification(notification)

	expectedSubmitJob := model.SubmitJob{
		Name:      "testJob",
		StartTime: 2,
		Interval:  1000,
		URL:       "http:localhost:6060/ready",
	}

	if actualSubmitJob != expectedSubmitJob {
		t.Errorf("Expected actual SubmitJob:\n'%v'to equal:\n'%v'\n", actualSubmitJob, expectedSubmitJob)
	}
}

type dummyLogger struct {
	logger    *log.Logger
	logBuffer bytes.Buffer
}

func newDummyLogger() *dummyLogger {
	dumLogger := dummyLogger{}

	dumLogger.logger = log.New(nil, "", 0)
	dumLogger.logger.SetOutput(&dumLogger.logBuffer)

	return &dumLogger
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
