package payload_test

import (
	"encoding/json"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/capture/payload"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayloadMarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payl     payload.Payload
		expected []byte
	}{
		{
			"with number",
			payload.Payload{
				payload.Metric{
					Name:  "temp",
					Value: 10,
				},
			},
			[]byte(`[{"name":"temp","value":10}]`),
		},
		{
			"with array",
			payload.Payload{
				payload.Metric{
					Name:  "power",
					Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0},
				},
			},
			[]byte(`[{"name":"power","value":[-78.75,-80.5,-73.75,-70.75,-72]}]`),
		},
		{
			"with  2 metrics as array",
			payload.Payload{
				payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
				payload.Metric{Name: "frequencies", Value: []interface{}{100.0, 200.0, 300.0, 400.0, 500.0}},
			},
			[]byte(`[{"name":"power","value":[-78.75,-80.5,-73.75,-70.75,-72]},{"name":"frequencies","value":[100,200,300,400,500]}]`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := json.Marshal(&tc.payl)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
