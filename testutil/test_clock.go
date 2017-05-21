package testutil

import "github.com/tjheslin1/GoSchedule/clock"

// TestClock contains an embedded `Clock` to override the functionality of
// the time.Time API function calls.
type TestClock struct {
	Clck *clock.Clock
}

// TimeNowUnixNano returns the initial time set in the `Clock`,
// plus 1000 milliseconds.
// It does not update the state of the underlying `Clock`.
func (testClock *TestClock) TimeNowUnixNano() int64 {
	return testClock.Clck.Now.UnixNano() + 1000
}
