package database

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/tjheslin1/GoSchedule/model"
	"github.com/tjheslin1/GoSchedule/testutil"
)

func TestPassNotification_SubmittedJob(t *testing.T) {
	notification := pq.Notification{BePid: -1, Channel: "", Extra: exampleNotification}

	submittedJobs := make(chan model.SubmitJob, 1)
	testLogger := testutil.NewTestLogger()

	jobListener := JobListener{submittedJobs, 1, testLogger.Logger}

	passNotitification(&notification, &jobListener)

	actualJob := <-submittedJobs
	if actualJob != expectedSubmitJob {
		t.Errorf("Expected submittedJob:\n%v\nbut got:\n%v\n", expectedSubmitJob, actualJob)
	}
}

// TODO callback is called many times, some with err populated.
func TestPassNotification_TimeAfter(t *testing.T) {
	t.SkipNow()
	submittedJobs := make(chan model.SubmitJob, 1)
	testLogger := testutil.NewTestLogger()

	jobListener := JobListener{submittedJobs, 1, testLogger.Logger}

	callback := func(eventType pq.ListenerEventType, err error) {
		if err != nil {
			t.Error("Expected callback error to be nil.")
		}
		fmt.Println("Callback Hit!")
	}

	jobListener.waitForNotification(pq.NewListener("name", 0, 0, callback))

	<-time.After(1 * time.Second)
	// if actualJob != expectedSubmitJob {
	// 	t.Errorf("Expected submittedJob:\n%v\nbut got:\n%v\n", expectedSubmitJob, actualJob)
	// }
}

func TestReportProblemCallback(t *testing.T) {
	testLogger := testutil.NewTestLogger()
	jobListener := JobListener{nil, 0, testLogger.Logger}

	jobListener.reportProblemCallback(0, errors.New("Test error"))

	var expectedLogOuput = "Error occured for pq.ListenerEventType: '0'.\nTest error\n"
	if testLogger.LogOutput() != expectedLogOuput {
		t.Errorf("Expected log output of:'%s'\nbut got:\n'%v'\n", expectedLogOuput, testLogger.LogOutput())
	}
}

func TestUnmarshallJSONJobNotification(t *testing.T) {
	testLogger := testutil.NewTestLogger()

	jobListener := JobListener{nil, 0, testLogger.Logger}
	actualSubmitJob := jobListener.unmarshallJSONJobNotification([]byte(exampleNotification))

	if actualSubmitJob != expectedSubmitJob {
		t.Errorf("Expected actual SubmitJob:\n'%v'to equal:\n'%v'\n", actualSubmitJob, expectedSubmitJob)
	}
}

const exampleNotification string = `{
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
	}`

var expectedSubmitJob = model.SubmitJob{
	Name:      "testJob",
	StartTime: 2,
	Interval:  1000,
	URL:       "http:localhost:6060/ready",
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
