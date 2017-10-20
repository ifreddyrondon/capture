package capture

import (
	"github.com/ifreddyrondon/gocapture/capture/geocoding"
)

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	*geocoding.Point
	*Date
	Payload interface{}
}

// NewCapture returns a new pointer to a Capture composed of the passed Point, Time and payload
func NewCapture(point *geocoding.Point, timestamp *Date, payload interface{}) *Capture {
	return &Capture{point, timestamp, payload}
}

// UnmarshalJSON decodes the capture from a JSON body.
// Throws an error if the body of the date cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (c *Capture) UnmarshalJSON(data []byte) error {
	p := geocoding.Point{}
	if err := p.UnmarshalJSON(data); err != nil {
		return err
	}

	date := Date{}
	date.UnmarshalJSON(data)

	capture := NewCapture(&p, &date, "")
	*c = *capture
	return nil
}
