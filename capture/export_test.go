package capture

// SetClockInstance is a helper function only exported for test.
// It's intended to be used for stub the time.Now() function.
func SetClockInstance(targetTimestamp *Timestamp, clock *Clock) {
	targetTimestamp.clock = clock
}
