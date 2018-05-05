package capture_test

import (
	"testing"
	"time"

	"gopkg.in/src-d/go-kallax.v1"

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
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{1, 2, 3}},
			},
			time.Now(),
			geocoding.Point{},
		},
		{
			"capture with point",
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{1, 2, 3}},
			},
			time.Now(),
			*point,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := capture.Capture{
				Payload:   tc.payload,
				Timestamp: tc.timestamp,
				Point:     tc.point,
			}
			require.NotNil(t, result)
			require.Equal(t, tc.payload, result.Payload)
			require.Equal(t, tc.point, result.Point)
			require.NotNil(t, result.Timestamp)
		})
	}
}

func TestCaptureUnmarshalWithOnlyPayload(t *testing.T) {
	t.Parallel()

	payloadData := payload.Payload{
		&payload.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
	}
	expected := capture.Capture{
		Payload:   payloadData,
		Timestamp: time.Now(),
		Point:     geocoding.Point{},
	}

	result := capture.Capture{}
	err := result.UnmarshalJSON([]byte(`{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}]}`))
	require.Nil(t, err)
	assert.Nil(t, result.LAT)
	assert.Nil(t, result.LNG)
	assert.NotNil(t, result.Timestamp)
	assert.NotNil(t, result.Tags)
	assert.Equal(t, expected.Payload, result.Payload)
}

func TestCaptureUnmarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	payl := payload.Payload{
		&payload.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
	}
	tt := []struct {
		name    string
		payload []byte
		result  *capture.Capture
	}{
		{
			"capture with payload timestamp",
			[]byte(`{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "date": "1989-12-26T06:01:00.00Z"}`),
			getCaptureWithoutPoint(payl, "1989-12-26T06:01:00.00Z"),
		},
		{
			"success with payload timestamp and point",
			[]byte(`{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 1),
		},
		{
			"success with payload timestamp and point with elevation",
			[]byte(`{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "lat": 1, "lng": 1, "elevation": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			getCaptureWithElevation(payl, "1989-12-26T06:01:00.00Z", 1, 1, 1),
		},
		{
			"success with payload timestamp and tags",
			[]byte(`{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "date": "1989-12-26T06:01:00.00Z", "tags": ["tag1", "tag2"]}`),
			getCaptureWithTags(payl, "1989-12-26T06:01:00.00Z", "tag1", "tag2"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := capture.Capture{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Equal(t, tc.result.LAT, result.LAT)
			assert.Equal(t, tc.result.LNG, result.LNG)
			assert.Equal(t, tc.result.Elevation, result.Elevation)
			assert.Equal(t, tc.result.Timestamp, result.Timestamp)
			assert.Equal(t, tc.result.Payload, result.Payload)
			assert.Equal(t, tc.result.Tags, result.Tags)
		})
	}
}

func TestCaptureUnmarshalJSONFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		payload []byte
		errs    []string
	}{
		{
			"invalid point",
			[]byte(`{"lat": -91, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			[]string{"latitude out of boundaries, may range from -90.0 to 90.0"},
		},
		{
			"missing point lat",
			[]byte(`{"lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			[]string{"latitude must not be blank"},
		},
		{
			"missing payload",
			[]byte(`{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			[]string{"payload value must not be blank"},
		},
		{
			"bad payload",
			[]byte(`{`),
			[]string{"cannot unmarshal json into valid capture"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := capture.Capture{}
			err := c.UnmarshalJSON(tc.payload)
			assert.Error(t, err)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

func TestCaptureMarshalJSON(t *testing.T) {
	t.Parallel()

	payl := payload.Payload{
		&payload.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
	}
	date := "1989-12-26T06:01:00.00Z"
	tt := []struct {
		name     string
		capture  *capture.Capture
		expected string
	}{
		{
			"capture with point",
			getCapture(payl, date, 1, 2),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2,"elevation":null}`,
		},
		{
			"capture with point and elevation",
			getCaptureWithElevation(payl, date, 1, 2, 3),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2,"elevation":3}`,
		},
		{
			"capture without a point",
			getCaptureWithoutPoint(payl, date),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":null,"lng":null,"elevation":null}`,
		},
		{
			"capture with tags",
			getCaptureWithTags(payl, date, "tag1", "tag2"),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":["tag1","tag2"],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":null,"lng":null,"elevation":null}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.capture
			// override auto generated fields for test purpose
			c.ID, _ = kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			c.CreatedAt, c.UpdatedAt = getDate(date), getDate(date)
			result, err := c.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}

func getCapture(p payload.Payload, date string, lat, lng float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	ts := getDate(date)
	return &capture.Capture{Payload: p, Timestamp: ts, Point: *point, Tags: []string{}}
}

func getCaptureWithElevation(p payload.Payload, date string, lat, lng, elevation float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	point.Elevation = &elevation
	ts := getDate(date)
	return &capture.Capture{Payload: p, Timestamp: ts, Point: *point, Tags: []string{}}
}

func getCaptureWithoutPoint(p payload.Payload, date string) *capture.Capture {
	ts := getDate(date)
	return &capture.Capture{Payload: p, Timestamp: ts, Point: geocoding.Point{}, Tags: []string{}}
}

func getCaptureWithTags(p payload.Payload, date string, tags ...string) *capture.Capture {
	ts := getDate(date)
	return &capture.Capture{Payload: p, Timestamp: ts, Point: geocoding.Point{}, Tags: tags}
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
