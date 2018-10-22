package decoder_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/capture/decoder"
	"github.com/ifreddyrondon/capture/features/capture/geocoding"
	locationDecoder "github.com/ifreddyrondon/capture/features/capture/geocoding/decoder"
	"github.com/ifreddyrondon/capture/features/capture/payload"
	payloadDecoder "github.com/ifreddyrondon/capture/features/capture/payload/decoder"
	tagsDecoder "github.com/ifreddyrondon/capture/features/capture/tags/decoder"
	timestampDecoder "github.com/ifreddyrondon/capture/features/capture/timestamp/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func s2n(v string) *json.Number {
	n := json.Number(v)
	return &n
}

func f2P(v float64) *float64 {
	return &v
}

func TestDecodePOSTCaptureOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     string
		expected decoder.POSTCapture
	}{
		{
			name: "decode capture with just payload",
			body: `{"payload":[{"name":"power","value":10}]}`,
			expected: decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
			},
		},
		{
			name: "decode capture with payload and location with lat lng",
			body: `{"payload":[{"name":"power","value":10}],"location":{"lat":1,"lng":1}}`,
			expected: decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				Location: &locationDecoder.PostPoint{LAT: f2P(1), LNG: f2P(1)},
			},
		},
		{
			name: "decode capture with payload and location with lat, lng and elevation",
			body: `{"payload":[{"name":"power","value":10}],"location":{"lat":1,"lng":1,"elevation":1}}`,
			expected: decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				Location: &locationDecoder.PostPoint{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
		{
			name: "decode capture with payload and timestamp",
			body: `{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z"}`,
			expected: decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				PostTimestamp: timestampDecoder.PostTimestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
			},
		},
		{
			name: "decode capture with payload, timestamp and location with lat, lng and elevation",
			body: `{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z","location":{"lat":1,"lng":1,"elevation":1}}`,
			expected: decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				PostTimestamp: timestampDecoder.PostTimestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
				Location:      &locationDecoder.PostPoint{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
		{
			name: "decode capture with payload, timestamp, location (lat, lng and elevation) and tags",
			body: `{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z","location":{"lat":1,"lng":1,"elevation":1},"tags":["at night"]}`,
			expected: decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				PostTimestamp: timestampDecoder.PostTimestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
				PostTags:      tagsDecoder.PostTags{Tags: []string{"at night"}},
				Location:      &locationDecoder.PostPoint{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var POSTCapture decoder.POSTCapture
			err := decoder.Decode(r, &POSTCapture)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.PostPayload, POSTCapture.PostPayload)
			assert.Equal(t, tc.expected.Location, POSTCapture.Location)
			assert.Equal(t, tc.expected.PostTimestamp.Timestamp, POSTCapture.PostTimestamp.Timestamp)
			assert.Equal(t, tc.expected.PostTimestamp.Date, POSTCapture.PostTimestamp.Date)
			assert.Equal(t, tc.expected.PostTags, POSTCapture.PostTags)
		})
	}
}

func TestDecodePOSTCaptureError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			"decode POSTCapture when invalid json",
			".",
			"cannot unmarshal json into valid capture",
		},
		{
			"decode POSTCapture when missing payload",
			`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`,
			"payload value must not be blank",
		},
		{
			"decode POSTCapture when invalid point",
			`{"payload":[{"name":"power","value":10}],"location":{"lat":-91,"lng":1}}`,
			"latitude out of boundaries, may range from -90.0 to 90.0",
		},
		{
			"decode POSTCapture when missing point lat",
			`{"payload":[{"name":"power","value":10}],"location":{"lng":1}}`,
			"latitude must not be blank",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var capt decoder.POSTCapture
			err := decoder.Decode(r, &capt)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestCaptureFromPOSTCaptureOK(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name        string
		POSTCapture decoder.POSTCapture
		expected    features.Capture
	}{
		{
			"get features.Capture from decoder.POSTCapture when point present and not defaults",
			decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				PostTimestamp: timestampDecoder.PostTimestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
				PostTags:      tagsDecoder.PostTags{Tags: []string{"at night"}},
				Location:      &locationDecoder.PostPoint{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
			features.Capture{
				Payload: payload.Payload{
					payload.Metric{Name: "power", Value: 10.0},
				},
				Tags:     []string{"at night"},
				Location: &geocoding.Point{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
		{
			"get features.Capture from decoder.POSTCapture without point and not defaults",
			decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
				PostTimestamp: timestampDecoder.PostTimestamp{Date: s2n("1989-12-26T06:01:00.00Z")},
				PostTags:      tagsDecoder.PostTags{Tags: []string{"at night"}},
			},
			features.Capture{
				Payload: payload.Payload{
					payload.Metric{Name: "power", Value: 10.0},
				},
				Tags: []string{"at night"},
			},
		},
		{
			"get features.Capture from decoder.POSTCapture without point and default time and tags",
			decoder.POSTCapture{
				PostPayload: payloadDecoder.PostPayload{
					Payload: payload.Payload{
						payload.Metric{Name: "power", Value: 10.0},
					},
				},
			},
			features.Capture{
				Payload: payload.Payload{
					payload.Metric{Name: "power", Value: 10.0},
				},
				Tags: []string{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			capt := tc.POSTCapture.GetCapture()
			assert.Equal(t, tc.expected.Payload, capt.Payload)
			assert.NotNil(t, capt.Timestamp)
			assert.Equal(t, tc.expected.Tags, capt.Tags)
			assert.Equal(t, tc.expected.Location, capt.Location)
			// test capture fields filled with not default values
			assert.NotEqual(t, kallax.ULID{}, capt.ID)
			assert.NotEqual(t, time.Time{}, capt.CreatedAt)
			assert.NotEqual(t, time.Time{}, capt.UpdatedAt)
		})
	}
}
