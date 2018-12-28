package adding_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func s2n(v string) *json.Number {
	n := json.Number(v)
	return &n
}

func TestValidateCaptureOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     string
		expected adding.Capture
	}{
		{
			name: "decode capture with just payload",
			body: `{"payload":[{"name":"power","value":10}]}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: 10.0},
					},
				},
			},
		},
		{
			name: "decode capture with payload and location with lat lng",
			body: `{"payload":[{"name":"power","value":10}],"location":{"lat":1,"lng":1}}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: 10.0},
					},
				},
				Location: &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1)},
			},
		},
		{
			name: "decode capture with payload and location with lat, lng and elevation",
			body: `{"payload":[{"name":"power","value":10}],"location":{"lat":1,"lng":1,"elevation":1}}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{{Name: "power", Value: 10.0}},
				},
				Location: &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
		{
			name: "decode capture with payload and timestamp",
			body: `{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z"}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{{Name: "power", Value: 10.0}},
				},
				Timestamp: adding.Timestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
			},
		},
		{
			name: "decode capture with payload, timestamp and location with lat, lng and elevation",
			body: `{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z","location":{"lat":1,"lng":1,"elevation":1}}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{{Name: "power", Value: 10.0}},
				},
				Timestamp: adding.Timestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
				Location:  &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
		{
			name: "decode capture with payload, timestamp, location (lat, lng and elevation) and tags",
			body: `{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z","location":{"lat":1,"lng":1,"elevation":1},"tags":["at night"]}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{{Name: "power", Value: 10.0}},
				},
				Timestamp: adding.Timestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
				Tags:      []string{"at night"},
				Location:  &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
		{
			name: "decode capture with payload value as array",
			body: `{"payload":[{"name":"power","value":[10, -1, -3]}]}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{{Name: "power", Value: []interface{}{10.0, -1.0, -3.0}}},
				},
			},
		},
		{
			name: "decode capture with several metrict and value as array",
			body: `{"payload":[{"name":"power","value":[10, -1, -3]},{"name":"signal","value":100}]}`,
			expected: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: []interface{}{10.0, -1.0, -3.0}},
						{Name: "signal", Value: 100.0},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var capture adding.Capture
			err := adding.CaptureValidator.Decode(r, &capture)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Payload, capture.Payload)
			assert.Equal(t, tc.expected.Location, capture.Location)
			assert.Equal(t, tc.expected.Timestamp.Timestamp, capture.Timestamp.Timestamp)
			assert.Equal(t, tc.expected.Timestamp.Date, capture.Timestamp.Date)
			assert.Equal(t, tc.expected.Tags, capture.Tags)
		})
	}
}

func TestValidationCaptureFails(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		errs []string
	}{
		{
			"decode Capture when invalid json",
			".",
			[]string{"cannot unmarshal json into valid capture value"},
		},
		{
			"decode Capture when missing payload",
			`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`,
			[]string{"payload value must not be blank"},
		},
		{
			"decode Capture when invalid point",
			`{"payload":[{"name":"power","value":10}],"location":{"lat":-91,"lng":1}}`,
			[]string{"latitude out of boundaries, may range from -90.0 to 90.0"},
		},
		{
			"decode Capture when missing point lat",
			`{"payload":[{"name":"power","value":10}],"location":{"lng":1}}`,
			[]string{"latitude must not be blank"},
		},
		{
			"decode Capture when invalid timestamp",
			`{"payload":[{"name":"power","value":10}],"timestamp":"a"}`,
			[]string{"invalid timestamp value: Could not find date format for a"},
		},
		{
			"decode Capture when invalid timestamp and location",
			`{"payload":[{"name":"power","value":10}],"timestamp":"a","location":{"lng":1}}`,
			[]string{
				"invalid timestamp value: Could not find date format for a",
				"latitude must not be blank",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var capture adding.Capture
			err := adding.CaptureValidator.Decode(r, &capture)
			for _, e := range tc.errs {
				assert.Contains(t, err.Error(), e)
			}
		})
	}
}
