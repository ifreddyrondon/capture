package captures

import (
	"time"

	"github.com/ifreddyrondon/gocapture/geocoding"
)

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	Point     *geocoding.Point
	Timestamp time.Time
	Payload   interface{}
}

// NewCapture returns a new pointer to a Capture composed of the passed Point, Time and payload
func NewCapture(point *geocoding.Point, timestamp time.Time, payload interface{}) *Capture {
	return &Capture{point, timestamp, payload}
}
