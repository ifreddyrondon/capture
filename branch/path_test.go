package branch_test

import (
	"testing"

	"time"

	"encoding/json"

	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyBranch(t *testing.T) {
	p := []byte(`[]`)

	var b branch.Branch
	err := b.UnmarshalJSON(p)
	require.Nil(t, err)
	require.Empty(t, b, "Expected len of branch to be 0. Got '%v'", len(b))
}

func TestPathUnmarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		payload []byte
		result  branch.Branch
	}{
		{
			"path of len 1",
			[]byte(`[{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.Branch{getCapture(1, 1, "1989-12-26T06:01:00.00Z", nil)},
		},
		{
			"path of len 2",
			[]byte(`[
			{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
			{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.Branch{
				getCapture(1, 1, "1989-12-26T06:01:00.00Z", nil),
				getCapture(1, 2, "1989-12-26T06:01:00.00Z", nil),
			},
		},
		{
			"invalid capture into path of len 1",
			[]byte(`[{"lat": -101, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.Branch{},
		},
		{
			"invalid capture into path of len 2",
			[]byte(`[
			{"lat": -101, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
			{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.Branch{getCapture(1, 2, "1989-12-26T06:01:00.00Z", nil)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b branch.Branch
			err := b.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Len(t, b, len(tc.result))
		})
	}
}

func TestMarshalBranch(t *testing.T) {
	date := "1989-12-26T06:01:00.00Z"

	c1 := getCapture(1, 2, date, []float64{12, 11})
	// override auto generated fields for test purpose
	c1.ID = "1" // the unmarshal of BsonId is an hexadecimal representation, e.g. "1"->"31"
	c1.CreatedDate, c1.LastModified = getDate(date), getDate(date)
	c2 := getCapture(5, 6, date, []float64{1, 2})
	c2.ID = "1"
	c2.CreatedDate, c2.LastModified = getDate(date), getDate(date)

	p := branch.Branch{c1, c2}
	result, _ := json.Marshal(p)
	expected := `[{"id":"31","payload":[12,11],"created_date":"1989-12-26T06:01:00Z","last_modified":"1989-12-26T06:01:00Z","timestamp":"1989-12-26T06:01:00Z","lat":1,"lng":2},{"id":"31","payload":[1,2],"created_date":"1989-12-26T06:01:00Z","last_modified":"1989-12-26T06:01:00Z","timestamp":"1989-12-26T06:01:00Z","lat":5,"lng":6}]`

	assert.Equal(t, expected, string(result))
}

func getCapture(lat, lng float64, date string, p []float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	payload := numberlist.New(p...)

	return capture.New(point, getDate(date), payload)
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
