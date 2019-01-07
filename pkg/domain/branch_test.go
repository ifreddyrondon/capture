package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg/domain"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSONBranch(t *testing.T) {
	t.Parallel()

	date := "1989-12-26T06:01:00.00Z"
	capturePayload := domain.Payload{
		domain.Metric{
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
	b := domain.Branch{
		ID:       bID,
		Name:     "master",
		Captures: []domain.Capture{c1, c2},
	}
	result, _ := json.Marshal(b)
	expected := `{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a395","name":"master","captures":[{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","repoId":"00000000-0000-0000-0000-000000000000"},{"id":"0162eb39-bd52-085b-3f0c-be3418244ec3","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","repoId":"00000000-0000-0000-0000-000000000000"}]}`

	assert.Equal(t, expected, string(result))
}

func getCapture(p domain.Payload, date string, lat, lng float64) domain.Capture {
	point := domain.Point{LAT: &lat, LNG: &lng}
	ts := getDate(date)
	return domain.Capture{Payload: p, Timestamp: ts, Location: &point, Tags: []string{}}
}
func getCaptureWithElevation(p domain.Payload, date string, lat, lng, elevation float64) domain.Capture {
	point := domain.Point{LAT: &lat, LNG: &lng, Elevation: &elevation}
	ts := getDate(date)
	return domain.Capture{Payload: p, Timestamp: ts, Location: &point, Tags: []string{}}
}
func getCaptureWithoutPoint(p domain.Payload, date string) domain.Capture {
	ts := getDate(date)
	return domain.Capture{Payload: p, Timestamp: ts, Tags: []string{}}
}
func getCaptureWithTags(p domain.Payload, date string, tags ...string) domain.Capture {
	ts := getDate(date)
	return domain.Capture{Payload: p, Timestamp: ts, Tags: tags}
}
func getDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
