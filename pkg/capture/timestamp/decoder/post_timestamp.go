package decoder

import (
	"encoding/json"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ifreddyrondon/capture/pkg"
)

type PostTimestamp struct {
	Date          *json.Number `json:"date"`
	Timestamp     *json.Number `json:"timestamp"`
	clock         *pkg.Clock
	postTimestamp *time.Time
}

func (t *PostTimestamp) OK() error {
	date := getNumber(t.Date, t.Timestamp)
	if date != nil {
		parsedTime, err := dateparse.ParseAny(date.String())
		if err != nil {
			return err
		}
		t.postTimestamp = &parsedTime
	}
	return nil
}

func (t *PostTimestamp) GetTimestamp() time.Time {
	if t.postTimestamp != nil {
		return t.postTimestamp.UTC()
	}
	return t.clock.Now()
}

func getNumber(data1, data2 *json.Number) *json.Number {
	if data1 == nil {
		return data2
	}
	return data1
}
