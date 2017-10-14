package captures

import (
	"errors"

	"github.com/ifreddyrondon/gocapture/geocoding"
)

var (
	CaptureUnmarshalError = errors.New("cannot unmarshal json into Capture value")
)

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	*geocoding.Point
	date    Date
	Payload interface{}
}

// NewCapture returns a new pointer to a Capture composed of the passed Point, Time and payload
func NewCapture(point *geocoding.Point, timestamp Date, payload interface{}) *Capture {
	return &Capture{point, timestamp, payload}
}
