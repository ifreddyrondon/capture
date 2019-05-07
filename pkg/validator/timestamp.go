package validator

import (
	"encoding/json"
	"time"

	"github.com/araddon/dateparse"
)

type Timestamp struct {
	Date      *json.Number `json:"date"`
	Timestamp *json.Number `json:"timestamp"`
	Time      *time.Time
}

func (t *Timestamp) Validate() error {
	date := getNumber(t.Date, t.Timestamp)
	if date != nil {
		parsedTime, err := dateparse.ParseAny(date.String())
		if err != nil {
			return err
		}
		t.Time = &parsedTime
	}
	return nil
}

func getNumber(data1, data2 *json.Number) *json.Number {
	if data1 == nil {
		return data2
	}
	return data1
}
