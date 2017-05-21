package model

import (
	"strings"
	"testing"
	"time"

	"github.com/tjheslin1/GoSchedule/clock"
	"github.com/tjheslin1/GoSchedule/testutil"
)

func TestStart(t *testing.T) {
	t.SkipNow()
	testLogger := testutil.NewTestLogger()
	testJob := SubmitJob{
		Name:      "TestJob",
		StartTime: 100,
		Interval:  2000,
		URL:       "http://test.com",
	}
	jobRun := JobRunner{
		Job:    testJob,
		Logger: testLogger.Logger,
	}

	startTime := time.Now()
	clck := clock.Clock{Now: startTime}
	testClock := testutil.TestClock{Clck: &clck}
	jobRun.Start(testClock.Clck)

	expectedLastFiredTime := startTime.UnixNano() + (1 * time.Second).Nanoseconds()
	if jobRun.lastFired != expectedLastFiredTime {
		t.Errorf("Expected lastFiredTime to be updated to:\n'%d'\nbut was:\n'%v'\n",
			expectedLastFiredTime, jobRun.lastFired)
	}

	expectedLogOutput := `Job Runner starting for '{TestJob 100 2000 http://test.com}'.
JOB FIRED!!!!
Job Runner finished for '{TestJob 100 2000 http://test.com}'.`

	if !strings.Contains(testLogger.LogOutput(), expectedLogOutput) {
		t.Errorf("Expected log output to be:\n%s\nbut was:\n%s\n", expectedLogOutput, testLogger.LogOutput())
	}
}

var testTimes = []struct {
	timeIn         time.Time
	startTime      int64
	interval       int64
	expectedResult bool
}{
	// year, month, day, hour, min, sec, nsec, loc
	{time.Date(1970, 1, 1, 0, 0, 1, 0, time.UTC), 2000, 200, false},
	{time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), 0, 50, true},
	{time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC), 0, 50, true},
	// TODO test interval triggers
	// {time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), 0, 50, true},
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
