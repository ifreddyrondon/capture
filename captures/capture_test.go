package captures_test

import (
	"os"
	"testing"
	"time"

	"github.com/ifreddyrondon/gocapture/captures"
	"github.com/ifreddyrondon/gocapture/geocoding"
)

func TestCaptureDateUnmarshalJSON(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	tt := []struct {
		name    string
		payload []byte
		result  string
	}{
		{
			"valid RFC3339 date with date key",
			[]byte(`{"date": "1989-12-26T06:01:00.00Z"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid RFC3339 date with timestamp key",
			[]byte(`{"timestamp": "1989-12-26T06:01:00.00Z"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid RFC1123 date with date key",
			[]byte(`{"date": "Tue, 26 Dec 1989 06:01:00 UTC"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid RFC1123 date with timestamp key",
			[]byte(`{"timestamp": "Tue, 26 Dec 1989 06:01:00 UTC"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid Unix timestamp date with date key",
			[]byte(`{"date": "630655260"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
		{
			"valid Unix timestamp date with timestamp key",
			[]byte(`{"timestamp": "630655260"}`),
			"1989-12-26 06:01:00 +0000 UTC",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := captures.NewCaptureDate()
			result.UnmarshalJSON(tc.payload)

			if result.Date.String() != tc.result {
				t.Errorf("Expected CaptureDate to be '%v'. Got '%v'", tc.result, result.Date.String())
			}
		})
	}
}

// TODO: add unit test when unmarshal capture date fails
//{"empty json", []byte("{}"), captures.CaptureDate{Date: nowFunc()}},

func TestNewCapture(t *testing.T) {
	point, _ := geocoding.NewPoint(1, 2)
	timestamp := captures.CaptureDate{Date: time.Now()}
	var payload interface{}

	result := captures.NewCapture(point, timestamp, payload)

	if result == nil {
		t.Errorf("Expected capture not to nil. Got '%v'", result)
	}
}
