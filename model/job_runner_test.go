package model

import (
	"testing"
	"time"

	"github.com/tjheslin1/GoSchedule/testutil"
)

var testTimes = []struct {
	timeIn         time.Time
	startTime      int64
	interval       int64
	expectedResult bool
}{
	{time.Date(1970, 1, 1, 0, 0, 1, 0, time.UTC), 2000, 200, false},
	{time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), 0, 50, true},
	{time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC), 0, 50, true},
}

func TestTimeToRun(t *testing.T) {
	for _, timeTest := range testTimes {

		testJob := SubmitJob{
			Name:      "TestJob",
			StartTime: timeTest.startTime,
			Interval:  timeTest.interval,
			URL:       "http://test.com",
		}
		testLogger := testutil.NewTestLogger()
		jobRun := JobRunner{
			Job:    testJob,
			Logger: testLogger.Logger,
		}

		isTimeToRun := jobRun.timeToRun(timeTest.timeIn)

		if isTimeToRun != timeTest.expectedResult {
			t.Errorf("Expected timeToRun to return '%v'. But got '%v'\nFor test input: '%v'\n",
				timeTest.expectedResult, isTimeToRun, timeTest)
		}
	}
}
