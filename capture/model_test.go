package capture_test

import (
	"testing"
	"time"

	"os"

	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
)

func TestNewCapture(t *testing.T) {
	point, _ := geocoding.NewPoint(1, 2)
	timestamp := capture.NewTimestamp(time.Now())
	var payload interface{}

	result := capture.NewCapture(point, timestamp, payload)

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
			getCapture(1, 1, "1989-12-26T06:01:00.00Z", ""),
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

			if result.Timestamp.Timestamp.String() != tc.result.Timestamp.Timestamp.String() {
				t.Errorf(
					"Expected Date of capture to be '%v'. Got '%v'",
					tc.result.Timestamp.Timestamp.String(),
					result.Timestamp.Timestamp.String())
			}
		})
	}
}

func getCapture(lat, lng float64, date string, payload interface{}) *capture.Capture {
	p, _ := geocoding.NewPoint(lat, lng)
	parsedDate, _ := time.Parse(time.RFC3339, date)
	timestamp := capture.NewTimestamp(parsedDate)

	return capture.NewCapture(p, timestamp, payload)
}