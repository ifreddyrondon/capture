package capture_test

import (
	"testing"
	"time"

	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/payload/numberlist"

	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCapture(t *testing.T) {
	p, err := geocoding.New(1, 2)
	payload := numberlist.New(12, 11)
	result := capture.New(p, time.Now(), payload)
	require.Nil(t, err)
	require.NotNil(t, result)
}

func TestCaptureUnmarshalJSONSuccess(t *testing.T) {
	tt := []struct {
		name    string
		payload []byte
		result  *capture.Capture
	}{
		{
			"success without payload",
			[]byte(`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			getCapture(1, 1, "1989-12-26T06:01:00.00Z", nil),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := capture.Capture{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Equal(t, tc.result.Lat, result.Lat)
			assert.Equal(t, tc.result.Lng, result.Lng)
			assert.Equal(t, tc.result.Timestamp, result.Timestamp)
			assert.Equal(t, tc.result.Payload, result.Payload)
		})
	}
}

func TestCaptureUnmarshalJSONFailure(t *testing.T) {
	tt := []struct {
		name    string
		payload []byte
		err     error
	}{
		{
			"invalid point",
			[]byte(`{"lat": -91, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			geocoding.ErrorLATRange,
		},
		{
			"missing point lat",
			[]byte(`{"lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			geocoding.ErrorLATMissing,
		},
		{
			"bad payload",
			[]byte(`{`),
			geocoding.ErrorUnmarshalPoint,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := capture.Capture{}
			err := c.UnmarshalJSON(tc.payload)
			require.Equal(t, tc.err, err)
		})
	}
}

func TestCaptureMarshalJSON(t *testing.T) {
	date := "1989-12-26T06:01:00.00Z"
	c := getCapture(1, 2, date, []float64{12, 11})
	// override auto generated fields for test purpose
	c.ID = "1" // the unmarshal of BsonId is an hexadecimal representation, e.g. "1"->"31"
	c.CreatedDate, c.LastModified = getDate(date), getDate(date)
	result, _ := c.MarshalJSON()
	expected := `{"id":"31","payload":[12,11],"created_date":"1989-12-26T06:01:00Z","last_modified":"1989-12-26T06:01:00Z","timestamp":"1989-12-26T06:01:00Z","lat":1,"lng":2}`

	assert.Equal(t, expected, string(result))
}

func getCapture(lat, lng float64, date string, p []float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	ts := getDate(date)
	payloadData := numberlist.New(p...)

	return capture.New(point, ts, payloadData)
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
