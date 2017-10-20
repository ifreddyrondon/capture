package capture

// SetClockInstance is a helper function only exported for test.
// It's intended to be used for stub the time.Now() function.
func SetClockInstance(targetDate *Date, clock *Clock) {
	targetDate.clock = clock
}
