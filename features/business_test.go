package features_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/capture/geocoding"
	"github.com/ifreddyrondon/capture/features/capture/payload"
	"github.com/ifreddyrondon/capture/features/user/decoder"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllowedVisibility(t *testing.T) {
	tt := []struct {
		name     string
		given    string
		expected bool
	}{
		{
			"empty visibility",
			"",
			false,
		},
		{
			"not allowed visibility",
			"protected",
			false,
		},
		{
			"public visibility",
			"public",
			true,
		},
		{
			"private visibility",
			"private",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := features.AllowedVisibility(tc.given)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMarshalJSONRepository(t *testing.T) {
	t.Parallel()

	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := []byte(`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","name":"test","current_branch":"","shared":true,"createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","owner":"00000000-0000-0000-0000-000000000000"}`)
	c := features.Repository{
		Name: "test",
		ID: func() kallax.ULID {
			id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			return id
		}(),
		Visibility: features.Public,
		CreatedAt:  d,
		UpdatedAt:  d,
	}

	result, err := json.Marshal(c)
	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

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
	b := features.Branch{
		ID:       bID,
		Name:     "master",
		Captures: []features.Capture{c1, c2},
	}
	result, _ := json.Marshal(b)
	expected := `{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a395","name":"master","captures":[{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"},{"id":"0162eb39-bd52-085b-3f0c-be3418244ec3","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}]}`

	assert.Equal(t, expected, string(result))
}

func TestMarshalJSONUser(t *testing.T) {
	t.Parallel()
	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := []byte(`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","email":"test@example.com","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`)
	u := features.User{
		Email: "test@example.com",
		ID: func() kallax.ULID {
			id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			return id
		}(),
		CreatedAt: d,
		UpdatedAt: d,
	}

	result, err := json.Marshal(u)
	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestCaptureMarshalJSON(t *testing.T) {
	t.Parallel()

	data := payload.Payload{
		payload.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
	}
	date := "1989-12-26T06:01:00.00Z"
	tt := []struct {
		name     string
		capture  features.Capture
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

func TestUserPassword(t *testing.T) {
	t.Parallel()

	email, password := "test@localhost.com", "b4KeHAYy3u9v=ZQX"
	postUser := decoder.PostUser{Email: &email, Password: &password}
	var u features.User
	err := postUser.User(&u)
	assert.Nil(t, err)
	assert.True(t, u.CheckPassword("b4KeHAYy3u9v=ZQX"))
	assert.False(t, u.CheckPassword("1"))
}

func getCapture(p payload.Payload, date string, lat, lng float64) features.Capture {
	point := geocoding.Point{LAT: &lat, LNG: &lng}
	ts := getDate(date)
	return features.Capture{Payload: p, Timestamp: ts, Location: &point, Tags: []string{}}
}
func getCaptureWithElevation(p payload.Payload, date string, lat, lng, elevation float64) features.Capture {
	point := geocoding.Point{LAT: &lat, LNG: &lng, Elevation: &elevation}
	ts := getDate(date)
	return features.Capture{Payload: p, Timestamp: ts, Location: &point, Tags: []string{}}
}
func getCaptureWithoutPoint(p payload.Payload, date string) features.Capture {
	ts := getDate(date)
	return features.Capture{Payload: p, Timestamp: ts, Tags: []string{}}
}
func getCaptureWithTags(p payload.Payload, date string, tags ...string) features.Capture {
	ts := getDate(date)
	return features.Capture{Payload: p, Timestamp: ts, Tags: tags}
}
func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
