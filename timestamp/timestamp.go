package timestamp

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"fmt"

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

	t.Timestamp = parsedTime
	return nil
}

// MarshalJSON decode current Date to JSON.
// It supports json.Marshaler interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", t.Timestamp.UTC())), nil
}
