package model

import (
	"log"
	"time"
)

// JobRunner runs in a go routine, checking whether or not it is time to execute
// the `Job` specified.
type JobRunner struct {
	Job    SubmitJob
	Logger *log.Logger
}

// Start begins the JobRunner which will exectute the `Job` at its next valid
// time. That is, at the next interval past the `Job's` `start_time`.
func (jRun *JobRunner) Start() {
	jRun.Logger.Println("Job Runner starting.")
	defer jRun.Logger.Println("Job Runner finished.")

	// for {
	// }
}

func (jRun *JobRunner) timeToRun(now time.Time) bool {
	millisNow := now.UnixNano() / 1000000

	if millisNow >= jRun.Job.StartTime {
		return true
	}

	return false
}
