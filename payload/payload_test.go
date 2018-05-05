package payload_test

import (
	"encoding/json"
	"testing"

	"github.com/ifreddyrondon/gocapture/payload"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayloadUnmarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payl     []byte
		expected payload.Payload
	}{
		{
			"with cap",
			[]byte(`{"cap": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}]}`),
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			},
		},
		{
			"with captures",
			[]byte(`{"captures": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}]}`),
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			},
		},
		{
			"with captures",
			[]byte(`{"data": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}]}`),
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			},
		},
		{
			"with captures",
			[]byte(`{"payload": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}]}`),
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			},
		},
		{
			"with 2 simples metrics ",
			[]byte(`{"payload": [{"name": "temp", "value": 10}, {"name": "power", "value": 30}]}`),
			payload.Payload{
				&payload.Metric{Name: "temp", Value: 10.0},
				&payload.Metric{Name: "power", Value: 30.0},
			},
		},
		{
			"with 2 metrics as array",
			[]byte(`{"payload": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}, {"name": "frequencies", "value": [100, 200, 300, 400, 500]}]}`),
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
				&payload.Metric{Name: "frequencies", Value: []interface{}{100.0, 200.0, 300.0, 400.0, 500.0}},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := payload.Payload{}
			err := result.UnmarshalJSON(tc.payl)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPayloadUnmarshalJSONFails(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		payload []byte
		errs    []string
	}{
		{
			"unmarshal error",
			[]byte(`'`),
			[]string{"cannot unmarshal json into valid payload value"},
		},
		{
			"unmarshal nil payload",
			[]byte(`{"payload": null`),
			[]string{"cannot unmarshal json into valid payload value"},
		},
		{
			"unmarshal empty payload",
			[]byte(`{"payload": []`),
			[]string{"cannot unmarshal json into valid payload value"},
		},
		{
			"unmarshal payload with nulls",
			[]byte(`{"payload": [null]`),
			[]string{"cannot unmarshal json into valid payload value"},
		},
		{
			"unmarshal empty body",
			[]byte(`{}`),
			[]string{"payload value must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := payload.Payload{}
			err := result.UnmarshalJSON(tc.payload)
			assert.Error(t, err)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

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
				&payload.Metric{
					Name:  "temp",
					Value: 10,
				},
			},
			[]byte(`[{"name":"temp","value":10}]`),
		},
		{
			"with array",
			payload.Payload{
				&payload.Metric{
					Name:  "power",
					Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0},
				},
			},
			[]byte(`[{"name":"power","value":[-78.75,-80.5,-73.75,-70.75,-72]}]`),
		},
		{
			"with  2 metrics as array",
			payload.Payload{
				&payload.Metric{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
				&payload.Metric{Name: "frequencies", Value: []interface{}{100.0, 200.0, 300.0, 400.0, 500.0}},
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
