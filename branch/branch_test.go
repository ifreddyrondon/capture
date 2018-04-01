package branch_test

import (
	"encoding/json"
	"testing"

	"time"

	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyBranch(t *testing.T) {
	t.Parallel()

	p := []byte(`[]`)
	var b branch.Branch
	err := b.UnmarshalJSON(p)
	require.Nil(t, err)
	require.Empty(t, b.Captures)
}

func TestPathUnmarshalJSON(t *testing.T) {
	t.Parallel()

	payl := map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}}
	tt := []struct {
		name    string
		payload []byte
		result  *branch.Branch
	}{
		{
			"path of len 1",
			[]byte(`[{"payload":{"power":[-70, -100.1, 3.1]}, "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.New("", getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 1)),
		},
		{
			"path of len 2",
			[]byte(`[
				{"payload":{"power":[-70, -100.1, 3.1]}, "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
				{"payload":{"power":[-70, -100.1, 3.1]}, "lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.New("", getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 1), getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 2)),
		},

		{
			"invalid capture (lat) into path of len 1",
			[]byte(`[{"lat": -101, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			&branch.Branch{},
		},
		{
			"invalid capture (missing payload) into path of len 1",
			[]byte(`[{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			&branch.Branch{},
		},
		{
			"invalid capture into path of len 2",
			[]byte(`[
				{"payload":{"power":[-70, -100.1, 3.1]}, "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
				{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			branch.New("", getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 2)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b branch.Branch
			err := b.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Len(t, b.Captures, len(tc.result.Captures))
		})
	}
}

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
