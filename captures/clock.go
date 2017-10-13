package captures

import "time"

// Clock interface is a drop replacement for time.Now().
// It's intended to be used as pointer fields on another struct.
// Using the ProductionClock that implements the Clock interface
// forward to the corresponding functions in the standard time package.
// You can use any struct that implements Clock to facilitates
// unit testing when accessing the current time.
type Clock interface {
	Now() time.Time
}

type ProductionClock struct{}

func (c *ProductionClock) Now() time.Time {
	return time.Now()
}
