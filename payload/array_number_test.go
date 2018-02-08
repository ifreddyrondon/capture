package payload_test

import (
	"testing"

	"encoding/json"

	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayNumberPayloadUnmarshalJSONSuccess(t *testing.T) {
	tt := []struct {
		name     string
		payload  []byte
		expected payload.ArrayNumberPayload
	}{
		{
			"unmarshal with cap",
			[]byte(`{"cap": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
		},
		{
			"unmarshal with captures",
			[]byte(`{"captures": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
		},
		{
			"unmarshal with data",
			[]byte(`{"data": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
		},
		{
			"unmarshal with payload",
			[]byte(`{"payload": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
		},
		{
			"unmarshal empty payload",
			[]byte(`{"payload": []}`),
			[]float64{},
		},
		{
			"unmarshal empty body",
			[]byte(`{}`),
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := payload.ArrayNumberPayload{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestArrayNumberPayloadUnmarshalJSON(t *testing.T) {
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := payload.ArrayNumberPayload{}
			err := result.UnmarshalJSON(tc.payload)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestArrayNumberPayloadMarshalJSON(t *testing.T) {
	tt := []struct {
		name     string
		payload  payload.ArrayNumberPayload
		expected string
	}{
		{"empty payload", payload.ArrayNumberPayload{}, `[]`},
		{
			"payload with valid data",
			payload.ArrayNumberPayload{1, 2, 3}, `[1,2,3]`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := json.Marshal(tc.payload)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
