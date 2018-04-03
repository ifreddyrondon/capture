package branch_test

import (
	"encoding/json"
	"testing"

	"time"

	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/stretchr/testify/assert"
)

func TestMarshalBranch(t *testing.T) {
	t.Parallel()

	date := "1989-12-26T06:01:00.00Z"
	payl := map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}}

	c1 := getCapture(payl, date, 1, 2)
	// override auto generated fields for test purpose
	c1.ID = 1 // the unmarshal of BsonId is an hexadecimal representation, e.g. "1"->"31"
	c1.CreatedAt, c1.UpdatedAt = getDate(date), getDate(date)
	c2 := getCapture(payl, date, 1, 2)
	c2.ID = 2
	c2.CreatedAt, c2.UpdatedAt = getDate(date), getDate(date)

	p := branch.New("", c1, c2)
	result, _ := json.Marshal(p)
	expected := `{"id":"","name":"master","captures":[{"id":1,"payload":{"power":[-70,-100.1,3.1]},"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2},{"id":2,"payload":{"power":[-70,-100.1,3.1]},"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2}]}`

	assert.Equal(t, expected, string(result))
}

func getCapture(p map[string]interface{}, date string, lat, lng float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	ts := getDate(date)
	capt, _ := capture.New(p, ts, *point)
	return capt
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
