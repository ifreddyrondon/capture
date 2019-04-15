package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

func TestMetricUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     []byte
		expected domain.Metric
	}{
		{
			"with array value",
			[]byte(`{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			domain.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
		},
		{
			"with literal numeric value",
			[]byte(`{"name": "power", "value": 123}`),
			domain.Metric{Name: "power", Value: 123.0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var result domain.Metric
			err := json.Unmarshal(tc.body, &result)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMetricMarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		metric   domain.Metric
		expected []byte
	}{
		{
			"with array value",
			domain.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			[]byte(`{"name":"power","value":[-78.75,-80.5,-73.75,-70.75,-72]}`),
		},
		{
			"with literal numeric value",
			domain.Metric{Name: "power", Value: 123.0},
			[]byte(`{"name":"power","value":123}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := json.Marshal(tc.metric)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCaptureMarshalJSON(t *testing.T) {
	t.Parallel()

	data := domain.Payload{
		domain.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
	}
	date := "1989-12-26T06:01:00.00Z"
	tt := []struct {
		name     string
		capture  domain.Capture
		expected string
	}{
		{
			"marshal capture with point",
			getCapture(data, date, 1, 2),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","repoId":"00000000-0000-0000-0000-000000000000"}`,
		},
		{
			"capture with point and elevation",
			getCaptureWithElevation(data, date, 1, 2, 3),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":{"lat":1,"lng":2,"elevation":3},"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","repoId":"00000000-0000-0000-0000-000000000000"}`,
		},
		{
			"capture without a point",
			getCaptureWithoutPoint(data, date),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":null,"tags":[],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","repoId":"00000000-0000-0000-0000-000000000000"}`,
		},
		{
			"capture with tags",
			getCaptureWithTags(data, date, "tag1", "tag2"),
			`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","payload":[{"name":"power","value":[-70,-100.1,3.1]}],"location":null,"tags":["tag1","tag2"],"timestamp":"1989-12-26T06:01:00Z","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","repoId":"00000000-0000-0000-0000-000000000000"}`,
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
