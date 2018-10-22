package features

import "time"

// Clock is a drop replacement for time.Now().
// It's intended to be used as pointer fields on another struct.
// Leaving the instance as a nil reference will cause any calls
// on the *Clock to forward to the corresponding functions in the
// standard time package. This is meant to be the behavior in production.
type Clock struct {
	instant time.Time
}

// Now returns the current time.Now() is the instance is nil or return the istance.
func (c *Clock) Now() time.Time {
	if c == nil {
		return time.Now()
	}
	return c.instant
}

// NewMockClock is a helper constructor to facilitates unit testing.
// It'll return the instant Time when Now() method is called.
func NewMockClock(instant time.Time) *Clock {
	return &Clock{instant: instant}
}
