package model

import (
	"log"
	"time"

	"github.com/tjheslin1/GoSchedule/clock"
)

// JobRunner runs in a go routine, checking whether or not it is time to execute
// the `Job` specified.
type JobRunner struct {
	Job       SubmitJob
	Logger    *log.Logger
	lastFired int64
}

// Start begins the JobRunner which will exectute the `Job` at its next valid
// time. That is, at the next interval past the `Job's` `start_time`.
func (job *JobRunner) Start(clck *clock.Clock) {
	job.Logger.Printf("Job Runner starting for '%v'.\n", job.Job)
	defer job.Logger.Printf("Job Runner finished for '%v'.\n", job.Job)

	for {
		if job.timeToRun(clck.Now) {
			job.Logger.Println("JOB FIRED!!!!")
			job.lastFired = clck.TimeNowUnixNano()

			// TODO continue to iterate.
			return
		}
		return
	}
}

// timeToRun compares the current Unix time to the `StartTime`
// specified in the `SubmitJob`.
func (job *JobRunner) timeToRun(now time.Time) bool {
	if unixNanoToEpoch(now.UnixNano()) >= job.Job.StartTime {
		return true
	}

	return false
}

// unixNanoToEpoch converts an int64 representing a time in nanoseconds,
// to the same time in milliseconds. This loses accuracy.
func unixNanoToEpoch(unixNano int64) int64 {
	return unixNano / 1000000
}
