package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
