package payload_test

import (
	json "encoding/json"
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
			"unmarshal with cap",
			[]byte(`{"cap": {"power": [-78.75, -80.5, -73.75, -70.75, -72]}}`),
			map[string]interface{}{"power": []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
		},
		{
			"unmarshal with captures",
			[]byte(`{"captures": {"power": [-78.75, -80.5, -73.75, -70.75, -72]}}`),
			map[string]interface{}{"power": []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
		},
		{
			"unmarshal with data",
			[]byte(`{"data": {"power": [-78.75, -80.5, -73.75, -70.75, -72]}}`),
			map[string]interface{}{"power": []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
		},
		{
			"unmarshal with payload",
			[]byte(`{"payload": {"power": [-78.75, -80.5, -73.75, -70.75, -72]}}`),
			map[string]interface{}{"power": []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
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

func TestArrayNumberPayloadUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payload  []byte
		expected error
	}{
		{
			"unmarshal error",
			[]byte(`'`),
			payload.ErrorUnmarshalPayload,
		},
		{
			"unmarshal empty payload",
			[]byte(`{"payload": {}}`),
			payload.ErrorMissingPayload,
		},
		{
			"unmarshal empty body",
			[]byte(`{}`),
			payload.ErrorMissingPayload,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := payload.Payload{}
			err := result.UnmarshalJSON(tc.payload)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestPayloadMarshalJSON(t *testing.T) {
	t.Parallel()

	expected := `{"power":[-78.75,-80.5,-73.75,-70.75,-72]}`
	payl := map[string]interface{}{"power": []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}}
	result, err := json.Marshal(payl)
	require.Nil(t, err)
	assert.Equal(t, expected, string(result))
}
