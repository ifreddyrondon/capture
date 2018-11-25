package decoder_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/capture/decoder"
	"github.com/ifreddyrondon/capture/pkg/capture/geocoding"
	locationDecoder "github.com/ifreddyrondon/capture/pkg/capture/geocoding/decoder"
	"github.com/ifreddyrondon/capture/pkg/capture/payload"
	payloadDecoder "github.com/ifreddyrondon/capture/pkg/capture/payload/decoder"
	tagsDecoder "github.com/ifreddyrondon/capture/pkg/capture/tags/decoder"
	timestampDecoder "github.com/ifreddyrondon/capture/pkg/capture/timestamp/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestDecodePUTCaptureOK(t *testing.T) {
	t.Parallel()

	id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")

	tt := []struct {
		name     string
		body     string
		expected decoder.PUTCapture
	}{
		{
			name: "decode capture with payload, timestamp, location (lat, lng and elevation) and tags",
			body: `{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z","location":{"lat":1,"lng":1,"elevation":1},"tags":["at night"]}`,
			expected: decoder.PUTCapture{
				ID: &id,
				POSTCapture: decoder.POSTCapture{
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
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var PUTCapture decoder.PUTCapture
			err := decoder.Decode(r, &PUTCapture)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.ID, PUTCapture.ID)
			assert.Equal(t, tc.expected.PostPayload, PUTCapture.PostPayload)
			assert.Equal(t, tc.expected.Location, PUTCapture.Location)
			assert.Equal(t, tc.expected.PostTimestamp.Timestamp, PUTCapture.PostTimestamp.Timestamp)
			assert.Equal(t, tc.expected.PostTimestamp.Date, PUTCapture.PostTimestamp.Date)
			assert.Equal(t, tc.expected.PostTags, PUTCapture.PostTags)
		})
	}
}

func TestDecodePUTCaptureError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			"decode PUTCapture when missing id",
			`{"payload":[{"name":"power","value":10}],"date":"1989-12-26T06:01:00.00Z","location":{"lat":1,"lng":1,"elevation":1},"tags":["at night"]}`,
			"capture id must not be blank",
		},
		{
			"decode PUTCapture when missing payload",
			`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`,
			"payload value must not be blank",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var capt decoder.PUTCapture
			err := decoder.Decode(r, &capt)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestCaptureFromPUTCaptureOK(t *testing.T) {
	t.Parallel()
	id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")

	tt := []struct {
		name       string
		PUTCapture decoder.PUTCapture
		expected   pkg.Capture
	}{
		{
			"get pkg.Capture from decoder.PUTCapture",
			decoder.PUTCapture{
				ID: &id,
				POSTCapture: decoder.POSTCapture{
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
			pkg.Capture{
				ID: id,
				Payload: payload.Payload{
					payload.Metric{Name: "power", Value: 10.0},
				},
				Tags:     []string{"at night"},
				Location: &geocoding.Point{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			capt := tc.PUTCapture.GetCapture()
			assert.Equal(t, tc.expected.Payload, capt.Payload)
			assert.NotNil(t, capt.Timestamp)
			assert.Equal(t, tc.expected.Tags, capt.Tags)
			assert.Equal(t, tc.expected.Location, capt.Location)
			assert.Equal(t, tc.expected.ID, capt.ID)
		})
	}
}
