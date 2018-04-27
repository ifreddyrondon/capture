package payload_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSONMetric(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payl     []byte
		expected payload.Metric
	}{
		{
			"with array value",
			[]byte(`{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
		},
		{
			"with literal numeric value",
			[]byte(`{"name": "power", "value": 123}`),
			payload.Metric{Name: "power", Value: 123.0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var result payload.Metric
			err := result.UnmarshalJSON(tc.payl)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMarshalJSONMetric(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		metric   payload.Metric
		expected []byte
	}{
		{
			"with array value",
			payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			[]byte(`{"name":"power","value":[-78.75,-80.5,-73.75,-70.75,-72]}`),
		},
		{
			"with literal numeric value",
			payload.Metric{Name: "power", Value: 123.0},
			[]byte(`{"name":"power","value":123}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.metric.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
