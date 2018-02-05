package timestamp

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/simplereach/timeutils"
)

// Timestamp represents the specific moment at which the capture was taken.
type Timestamp struct {
	Timestamp time.Time `json:"timestamp"`
	clock     *Clock
}

func NewTimestamp(date time.Time) *Timestamp {
	return &Timestamp{Timestamp: date}
}

type timestampJSON struct {
	stringer struct {
		Date      string `json:"date"`
		Timestamp string `json:"timestamp"`
	}
	Date      json.Number `json:"date"`
	Timestamp json.Number `json:"timestamp"`
}

// UnmarshalJSON decodes the Timestamp of the capture from a JSON body.
// Throws an error if the body of the Timestamp cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	t.Timestamp = t.clock.Now()

	model := new(timestampJSON)
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&model); err != nil {
		log.Print(err)
		return nil
	}

	var date string
	if model.Date != "" {
		date = model.Date.String()
	} else if model.Timestamp != "" {
		date = model.Timestamp.String()
	}

	parsedTime, err := timeutils.ParseDateString(date)
	if err != nil {
		return nil
	}

	t.Timestamp = parsedTime.UTC()
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if y := t.Timestamp.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = t.Timestamp.UTC().AppendFormat(b, time.RFC3339Nano)
	b = append(b, '"')
	return b, nil
}
