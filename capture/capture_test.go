package capture_test

import (
	"testing"
	"time"

	"github.com/ifreddyrondon/gocapture/payload"

	"github.com/ifreddyrondon/gocapture/geocoding"

	"github.com/ifreddyrondon/gocapture/capture"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCapture(t *testing.T) {
	t.Parallel()
	point, _ := geocoding.New(1, 2)
	tt := []struct {
		name      string
		payload   payload.Payload
		timestamp time.Time
		point     geocoding.Point
	}{
		{
			"simple capture: payload and timestamp",
			map[string]interface{}{"power": []float64{1, 2, 3}},
			time.Now(),
			geocoding.Point{},
		},
		{
			"capture with point",
			map[string]interface{}{"power": []float64{1, 2, 3}},
			time.Now(),
			*point,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := capture.New(tc.payload, tc.timestamp, tc.point)
			require.Nil(t, err)
			require.NotNil(t, result)
			require.Equal(t, tc.payload, result.Payload)
			require.Equal(t, tc.point, result.Point)
			require.NotNil(t, result.Timestamp)
		})
	}
}

func TestNewCaptureFail(t *testing.T) {
	t.Parallel()

	point, _ := geocoding.New(1, 2)
	_, err := capture.New(nil, time.Now(), *point)
	assert.EqualError(t, err, "missing payload value")
}

func TestCaptureUnmarshalWithOnlyPayload(t *testing.T) {
	t.Parallel()

	payloadData := map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}}
	expected, _ := capture.New(payloadData, time.Now(), geocoding.Point{})

	result := capture.Capture{}
	err := result.UnmarshalJSON([]byte(`{"payload":{"power":[-70, -100.1, 3.1]}}`))
	require.Nil(t, err)
	assert.Nil(t, result.LAT)
	assert.Nil(t, result.LNG)
	assert.NotNil(t, result.Timestamp)
	assert.Equal(t, expected.Payload, result.Payload)
}

func TestCaptureUnmarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	payl := map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}}
	tt := []struct {
		name    string
		payload []byte
		result  *capture.Capture
	}{
		{
			"capture with payload timestamp",
			[]byte(`{"payload":{"power":[-70, -100.1, 3.1]}, "date": "1989-12-26T06:01:00.00Z"}`),
			getCaptureWithoutPoint(payl, "1989-12-26T06:01:00.00Z"),
		},
		{
			"success with payload timestamp and point",
			[]byte(`{"payload":{"power":[-70, -100.1, 3.1]}, "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := capture.Capture{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Equal(t, tc.result.LAT, result.LAT)
			assert.Equal(t, tc.result.LNG, result.LNG)
			assert.Equal(t, tc.result.Timestamp, result.Timestamp)
			assert.Equal(t, tc.result.Payload, result.Payload)
		})
	}
}

func TestCaptureUnmarshalJSONFailure(t *testing.T) {
	t.Parallel()

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
			"missing payload",
			[]byte(`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			payload.ErrorMissingPayload,
		},
		{
			"bad payload",
			[]byte(`{`),
			capture.ErrorBadPayload,
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
	t.Parallel()

	payl := map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}}
	date := "1989-12-26T06:01:00.00Z"
	tt := []struct {
		name     string
		capture  *capture.Capture
		expected string
	}{
		{
			"capture with point",
			getCapture(payl, date, 1, 2),
			`{"id":1,"payload":{"power":[-70,-100.1,3.1]},"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2}`,
		},
		{
			"capture without a point",
			getCaptureWithoutPoint(payl, date),
			`{"id":1,"payload":{"power":[-70,-100.1,3.1]},"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":null,"lng":null}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.capture
			// override auto generated fields for test purpose
			c.ID = 1
			c.CreatedAt, c.UpdatedAt = getDate(date), getDate(date)
			result, _ := c.MarshalJSON()

			assert.Equal(t, tc.expected, string(result))
		})
	}
}

func getCapture(p map[string]interface{}, date string, lat, lng float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	ts := getDate(date)
	capt, _ := capture.New(p, ts, *point)
	return capt
}

func getCaptureWithoutPoint(p map[string]interface{}, date string) *capture.Capture {
	ts := getDate(date)
	capt, _ := capture.New(p, ts, geocoding.Point{})
	return capt
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
