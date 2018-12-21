package pkg_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/capture/geocoding"
	"github.com/ifreddyrondon/capture/pkg/capture/payload"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSONBranch(t *testing.T) {
	t.Parallel()

	date := "1989-12-26T06:01:00.00Z"
	capturePayload := payload.Payload{
		payload.Metric{
			Name:  "power",
			Value: []interface{}{-70.0, -100.1, 3.1},
		},
	}

	c1 := getCapture(capturePayload, date, 1, 2)
	// override auto generated fields for test purpose
	c1.ID, _ = kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
	c1.CreatedAt, c1.UpdatedAt = getDate(date), getDate(date)
	c2 := getCapture(capturePayload, date, 1, 2)
	c2.ID, _ = kallax.NewULIDFromText("0162eb39-bd52-085b-3f0c-be3418244ec3")
	c2.CreatedAt, c2.UpdatedAt = getDate(date), getDate(date)

	bID, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a395")
	b := pkg.Branch{
		ID:       bID,
		Name:     "master",
		Captures: []pkg.Capture{c1, c2},
	}
	result, _ := json.Marshal(b)
	expected := `{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a395","name":"master","captures":[{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"},{"id":"0162eb39-bd52-085b-3f0c-be3418244ec3","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}]}`

	assert.Equal(t, expected, string(result))
}

func TestCaptureMarshalJSON(t *testing.T) {
	t.Parallel()

	data := payload.Payload{
		payload.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
	}
	date := "1989-12-26T06:01:00.00Z"
	tt := []struct {
		name     string
		capture  pkg.Capture
		expected string
	}{
		{
			"marshal capture with point",
			getCapture(data, date, 1, 2),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`,
		},
		{
			"capture with point and elevation",
			getCaptureWithElevation(data, date, 1, 2, 3),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2,"elevation":3},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`,
		},
		{
			"capture without a point",
			getCaptureWithoutPoint(data, date),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":null,"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`,
		},
		{
			"capture with tags",
			getCaptureWithTags(data, date, "tag1", "tag2"),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":null,"tags":["tag1","tag2"],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.capture
			// override auto generated fields for test purpose
			c.ID, _ = kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			c.CreatedAt, c.UpdatedAt = getDate(date), getDate(date)
			result, err := json.Marshal(c)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}

func getCapture(p payload.Payload, date string, lat, lng float64) pkg.Capture {
	point := geocoding.Point{LAT: &lat, LNG: &lng}
	ts := getDate(date)
	return pkg.Capture{Payload: p, Timestamp: ts, Location: &point, Tags: []string{}}
}
func getCaptureWithElevation(p payload.Payload, date string, lat, lng, elevation float64) pkg.Capture {
	point := geocoding.Point{LAT: &lat, LNG: &lng, Elevation: &elevation}
	ts := getDate(date)
	return pkg.Capture{Payload: p, Timestamp: ts, Location: &point, Tags: []string{}}
}
func getCaptureWithoutPoint(p payload.Payload, date string) pkg.Capture {
	ts := getDate(date)
	return pkg.Capture{Payload: p, Timestamp: ts, Tags: []string{}}
}
func getCaptureWithTags(p payload.Payload, date string, tags ...string) pkg.Capture {
	ts := getDate(date)
	return pkg.Capture{Payload: p, Timestamp: ts, Tags: tags}
}
func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
