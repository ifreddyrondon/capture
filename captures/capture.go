package captures

import (
	"errors"
	"time"

	"bytes"
	"encoding/json"
	"log"

	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/simplereach/timeutils"
)

var (
	CaptureUnmarshalError = errors.New("cannot unmarshal json into Capture value")
)

var clockInstance Clock = new(ProductionClock)

type CaptureDate struct {
	Date  time.Time
	clock Clock
}

func NewCaptureDate() *CaptureDate {
	return &CaptureDate{clock: clockInstance}
}

// UnmarshalJSON decodes the date of the capture from a JSON body.
// Throws an error if the body of the date cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (t *CaptureDate) UnmarshalJSON(data []byte) error {
	t.Date = t.clock.Now()
	decoder := json.NewDecoder(bytes.NewReader(data))
	var values map[string]string
	if err := decoder.Decode(&values); err != nil {
		log.Print(err)
		return nil
	}

	var date string
	if val, isOk := getTimestampValue(&values); isOk {
		date = val
	}

	if date == "" {
		return nil
	}

	parsedTime, err := timeutils.ParseDateString(date)
	if err != nil {
		return nil
	}

	t.Date = parsedTime
	return nil
}

func getTimestampValue(values *map[string]string) (string, bool) {
	var val string
	var isOk bool
	if val, isOk = (*values)["date"]; isOk {
		isOk = true
	} else if val, isOk = (*values)["timestamp"]; isOk {
		isOk = true
	}

	return val, isOk
}

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	*geocoding.Point
	date    CaptureDate
	Payload interface{}
}

// NewCapture returns a new pointer to a Capture composed of the passed Point, Time and payload
func NewCapture(point *geocoding.Point, timestamp CaptureDate, payload interface{}) *Capture {
	return &Capture{point, timestamp, payload}
}
