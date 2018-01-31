package capture_test

import (
	"os"
	"testing"
	"time"

	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/ifreddyrondon/gocapture/timestamp"
)

func TestNewCapture(t *testing.T) {
	point, _ := geocoding.NewPoint(1, 2)
	ts := timestamp.NewTimestamp(time.Now())

	result := capture.NewCapture(point, ts, []float64{12, 11})

	if result == nil {
		t.Errorf("Expected capture not to nil. Got '%v'", result)
	}
}

func TestCaptureUnmarshalJSON(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	tt := []struct {
		name      string
		payload   []byte
		result    *capture.Capture
		resultErr error
	}{
		{
			"valid point with given date",
			[]byte(`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			getCapture(1, 1, "1989-12-26T06:01:00.00Z", []float64{}),
			nil,
		},
		{
			"invalid point",
			[]byte(`{"lat": -91, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			nil,
			geocoding.ErrorLATRange,
		},
		{
			"missing point lat",
			[]byte(`{"lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			nil,
			geocoding.ErrorLATMissing,
		},
		{
			"bad payload",
			[]byte(`{`),
			nil,
			geocoding.ErrorUnmarshalPoint,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := capture.Capture{}
			resultError := result.UnmarshalJSON(tc.payload)

			if resultError != tc.resultErr {
				t.Errorf("Expected get the error '%v'. Got '%v'", tc.resultErr, resultError)
			}

			// if result expected an error do not check for internal attrs
			if tc.resultErr != nil {
				return
			}

			if result.Point.Lat != tc.result.Point.Lat {
				t.Errorf("Expected Lat of capture to be '%v'. Got '%v'", tc.result.Point.Lat, result.Point.Lat)
			}

			if result.Point.Lng != tc.result.Point.Lng {
				t.Errorf("Expected Lng of capture to be '%v'. Got '%v'", tc.result.Point.Lng, result.Point.Lng)
			}

			if !result.Timestamp.Timestamp.Equal(tc.result.Timestamp.Timestamp) {
				t.Errorf(
					"Expected Date of capture to be '%v'. Got '%v'",
					tc.result.Timestamp.Timestamp, result.Timestamp.Timestamp)
			}

			if len(result.Payload) != len(tc.result.Payload) {
				t.Errorf("Expected payload to be '%v'. Got '%v'", len(tc.result.Payload), len(result.Payload))
			}

			for i, v := range tc.result.Payload {
				if v != result.Payload[i] {
					t.Fatalf("Expected payload at index %v to be '%v'. Got '%v'", i, v, result.Payload[i])
				}
			}
		})
	}
}

func getCapture(lat, lng float64, date string, payload payload.ArrayNumberPayload) *capture.Capture {
	p, _ := geocoding.NewPoint(lat, lng)
	parsedDate, _ := time.Parse(time.RFC3339, date)
	ts := timestamp.NewTimestamp(parsedDate)

	return capture.NewCapture(p, ts, payload)
}
