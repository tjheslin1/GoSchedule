package clock

import (
	"testing"
	"time"
)

func TestTimeNowUnixNano(t *testing.T) {
	now := time.Now()
	clock := Clock{Now: now}

	updatedNow := clock.TimeNowUnixNano()

	if clock.Now.UnixNano() != updatedNow {
		t.Errorf("Expected `Now` field of `Clock` to be updated to '%d' but was '%d'",
			updatedNow, clock.Now.UnixNano())
	}

	if clock.Now.UnixNano() <= now.UnixNano() {
		t.Errorf("Expected `Now` field of `Clock` to be updated and therefore "+
			"greater than '%d' but was '%d'",
			now.UnixNano(), clock.Now.UnixNano())
	}
}
