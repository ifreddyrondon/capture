package branch_test

import (
	"testing"

	"time"

	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/geocoding"
)

func TestPathAddCapture(t *testing.T) {
	tt := []struct {
		name     string
		captures []*capture.Capture
	}{
		{"empty path", []*capture.Capture{}},
		{"empty path", []*capture.Capture{
			getCapture(1, 1, "1989-12-26T06:01:00.00Z", ""),
			getCapture(1, 2, "1989-12-26T06:01:00.00Z", ""),
		}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := new(branch.Branch)
			p.AddCapture(tc.captures...)

			if p == nil {
				t.Errorf("Expected path not to nil. Got '%v'", p)
			}

			if len(p.Captures) != len(tc.captures) {
				t.Errorf("Expected path len to be '%v'. Got '%v'", len(tc.captures), len(p.Captures))
			}
		})
	}
}

func TestPathUnmarshalJSON(t *testing.T) {
	tt := []struct {
		name        string
		body        []byte
		resultError error
		expectedLen int
	}{
		{
			"empty path",
			[]byte(`[]`),
			nil,
			0,
		},
		{
			"valid capture into path of len 1",
			[]byte(`[{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			nil,
			1,
		},
		{
			"valid capture into path of len 2",
			[]byte(`[
				{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
				{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			nil,
			2,
		},
		{
			"invalid capture into path of len 1",
			[]byte(`[{"lat": -101, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
			nil,
			0,
		},
		{
			"invalid capture into path of len 2",
			[]byte(`[
				{"lat": -101, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
				{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
			nil,
			1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := new(branch.Branch)
			err := p.UnmarshalJSON(tc.body)

			if err != tc.resultError {
				t.Fatalf("Expected get the error '%v'. Got '%v'", tc.resultError, err)
			}

			if tc.resultError != nil {
				return
			}

			if tc.expectedLen != len(p.Captures) {
				t.Errorf("Expected len of catures to be '%v'. Got '%v'", tc.expectedLen, len(p.Captures))
			}
		})
	}
}

func getCapture(lat, lng float64, date string, payload interface{}) *capture.Capture {
	p, _ := geocoding.NewPoint(lat, lng)
	parsedDate, _ := time.Parse(time.RFC3339, date)
	timestamp := capture.NewDate(parsedDate)

	return capture.NewCapture(p, timestamp, payload)
}
