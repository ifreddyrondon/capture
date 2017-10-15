package captures

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/simplereach/timeutils"
)

// Date represents the specific timestamp at which the capture was taken.
type Date struct {
	Timestamp time.Time
	clock     *Clock
}

func NewDate(date time.Time) *Date {
	return &Date{Timestamp: date}
}

type dateJSON struct {
	stringer struct {
		Date      string `json:"date"`
		Timestamp string `json:"timestamp"`
	}
	Date      json.Number `json:"date"`
	Timestamp json.Number `json:"timestamp"`
}

// UnmarshalJSON decodes the date of the capture from a JSON body.
// Throws an error if the body of the date cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (t *Date) UnmarshalJSON(data []byte) error {
	t.Timestamp = t.clock.Now()

	model := new(dateJSON)
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
