package captures

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/simplereach/timeutils"
)

var clockInstance Clock = new(ProductionClock)

type Date struct {
	Timestamp time.Time
	clock     Clock
}

func NewDate() *Date {
	return &Date{clock: clockInstance}
}

// UnmarshalJSON decodes the date of the capture from a JSON body.
// Throws an error if the body of the date cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (t *Date) UnmarshalJSON(data []byte) error {
	t.Timestamp = t.clock.Now()
	decoder := json.NewDecoder(bytes.NewReader(data))
	var values map[string]interface{}
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

	t.Timestamp = parsedTime
	return nil
}

func getTimestampValue(values *map[string]interface{}) (string, bool) {
	var val interface{}
	var isOk bool
	if val, isOk = (*values)["date"]; isOk {
		isOk = true
	} else if val, isOk = (*values)["timestamp"]; isOk {
		isOk = true
	}

	if !isOk {
		return "", isOk
	}

	var date string
	switch v := val.(type) {
	case float64:
		date = strconv.Itoa(int(v))
	case string:
		date = v
	default:
		return "", false
	}

	return date, isOk
}
