package numberlist_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/payload"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	data := []float64{1, 2, 5.3, 4, 8}
	p := numberlist.New(data...)
	require.NotNil(t, p)
	require.Len(t, *p, len(data))
}

func TestArrayNumberPayloadUnmarshalJSONSuccess(t *testing.T) {
	tt := []struct {
		name     string
		payload  []byte
		expected numberlist.Payload
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
			result := numberlist.Payload{}
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
			numberlist.ErrorUnmarshalPayload,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := numberlist.Payload{}
			err := result.UnmarshalJSON(tc.payload)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestArrayNumberPayloadMarshalJSON(t *testing.T) {
	tt := []struct {
		name     string
		payload  numberlist.Payload
		expected string
	}{
		{"empty payload", numberlist.Payload{}, `[]`},
		{
			"payload with valid data",
			numberlist.Payload{1, 2, 3}, `[1,2,3]`,
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
