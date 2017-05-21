package clock

import "time"

// Clock wraps the provided time.Time for ease of testing,
// and to encapsulate useful time related functions.s
type Clock struct {
	Now time.Time
}

// TimeNowUnixNano updates `Now` field and returns the captured time
// as UnixNano.
func (clock *Clock) TimeNowUnixNano() int64 {
	clock.Now = time.Now()
	return clock.Now.UnixNano()
}
