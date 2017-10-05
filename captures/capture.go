package captures

import (
	"time"

	"github.com/ifreddyrondon/gocapture/geocoding"
)

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	Point     geocoding.Point
	Timestamp time.Duration
	Payload   interface{}
}
