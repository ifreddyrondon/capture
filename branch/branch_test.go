package branch_test

import (
	"encoding/json"
	"testing"

	"github.com/ifreddyrondon/gocapture/payload"

	"gopkg.in/src-d/go-kallax.v1"

	"time"

	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/stretchr/testify/assert"
)

func TestMarshalBranch(t *testing.T) {
	t.Parallel()

	date := "1989-12-26T06:01:00.00Z"
	payl := payload.Payload{
		&payload.Metric{
			Name:  "power",
			Value: []interface{}{-70.0, -100.1, 3.1},
		},
	}

	c1 := getCapture(payl, date, 1, 2)
	// override auto generated fields for test purpose
	c1.ID, _ = kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
	c1.CreatedAt, c1.UpdatedAt = getDate(date), getDate(date)
	c2 := getCapture(payl, date, 1, 2)
	c2.ID, _ = kallax.NewULIDFromText("0162eb39-bd52-085b-3f0c-be3418244ec3")
	c2.CreatedAt, c2.UpdatedAt = getDate(date), getDate(date)

	p := branch.New("", c1, c2)
	result, _ := json.Marshal(p)
	expected := `{"id":"","name":"master","captures":[{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":null,"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2,"elevation":null},{"id":"0162eb39-bd52-085b-3f0c-be3418244ec3","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":null,"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2,"elevation":null}]}`

	assert.Equal(t, expected, string(result))
}

func getCapture(p payload.Payload, date string, lat, lng float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	ts := getDate(date)
	return &capture.Capture{Payload: p, Timestamp: ts, Point: *point}
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
