package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/tjheslin1/GoSchedule/model"
	"github.com/tjheslin1/GoSchedule/testutil"
)

func TestReportProblemCallback(t *testing.T) {
	testLogger := testutil.NewTestLogger()

	jobListener := JobListener{testLogger.Logger}

	jobListener.reportProblemCallback(0, errors.New("Test error"))

	var expectedLogOuput = "Error occured for pq.ListenerEventType: '0'.\nTest error\n"
	if testLogger.LogOutput() != expectedLogOuput {
		t.Errorf("Expected log output of:'%s'\nbut got:\n'%v'\n", expectedLogOuput, testLogger.LogOutput())
	}
}

func TestUnmarshallJSONJobNotification(t *testing.T) {
	testLogger := testutil.NewTestLogger()

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

	jobListener := JobListener{testLogger.Logger}
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
