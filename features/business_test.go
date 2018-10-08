package features_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/capture"
	"github.com/ifreddyrondon/capture/features/capture/geocoding"
	"github.com/ifreddyrondon/capture/features/capture/payload"
	"github.com/ifreddyrondon/capture/features/user/decoder"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		Shared:    true,
		CreatedAt: d,
		UpdatedAt: d,
	}

	result, err := json.Marshal(c)
	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshalJSONBranch(t *testing.T) {
	t.Parallel()

	date := "1989-12-26T06:01:00.00Z"
	capturePayload := payload.Payload{
		&payload.Metric{
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
		Captures: []capture.Capture{*c1, *c2},
	}
	result, _ := json.Marshal(b)
	expected := `{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a395","name":"master","captures":[{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":null,"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2,"elevation":null},{"id":"0162eb39-bd52-085b-3f0c-be3418244ec3","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"tags":null,"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","lat":1,"lng":2,"elevation":null}]}`

	assert.Equal(t, expected, string(result))
}

func getCapture(p payload.Payload, date string, lat, lng float64) *capture.Capture {
	point, _ := geocoding.New(lat, lng)
	ts := getDate(date)
	return &capture.Capture{Payload: p, Timestamp: ts, Point: *point}
}

func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}

func TestUserPassword(t *testing.T) {
	t.Parallel()

	email, password := "test@localhost.com", "b4KeHAYy3u9v=ZQX"
	var u features.User
	err := decoder.User(decoder.PostUser{Email: &email, Password: &password}, &u)
	assert.Nil(t, err)
	assert.True(t, u.CheckPassword("b4KeHAYy3u9v=ZQX"))
	assert.False(t, u.CheckPassword("1"))
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
