package branch_test

import (
	"testing"

	"time"

	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/ifreddyrondon/gocapture/timestamp"
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
			getBranch(getCapture(1, 1, "1989-12-26T06:01:00.00Z", nil)),
		},
		{
			"path of len 2",
			[]byte(`[
			{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
			{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			getBranch(
				getCapture(1, 1, "1989-12-26T06:01:00.00Z", nil),
				getCapture(1, 2, "1989-12-26T06:01:00.00Z", nil),
			),
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
			getBranch(getCapture(1, 2, "1989-12-26T06:01:00.00Z", nil)),
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

func getBranch(captures ...*capture.Capture) branch.Branch {
	var b branch.Branch
	b = append(b, captures...)
	return b
}

func getCapture(lat, lng float64, date string, payload payload.ArrayNumberPayload) *capture.Capture {
	p, _ := geocoding.NewPoint(lat, lng)
	t, _ := time.Parse(time.RFC3339, date)
	ts := timestamp.NewTimestamp(t)

	return capture.NewCapture(p, ts, payload)
}
