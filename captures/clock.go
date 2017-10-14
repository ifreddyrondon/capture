package captures

import "time"

// Clock interface is a drop replacement for time.Now().
// It's intended to be used as pointer fields on another struct.
// You can use any struct that implements Clock to facilitates
// unit testing when accessing the current time.
type Clock interface {
	Now() time.Time
}

// ProductionClock implements the Clock interface and forward
// to the corresponding functions in the standard time package.
type ProductionClock struct{}

func (c *ProductionClock) Now() time.Time {
	return time.Now()
}

// MockClock implements the Clock interface and get an instant (time.Time)
// attribute. When Now() method is called, it returns the instant.
type MockClock struct {
	instant time.Time
}

func (c *MockClock) Now() time.Time {
	return c.instant
}

func NewMockClock(instant time.Time) *MockClock {
	return &MockClock{instant: instant}
}
