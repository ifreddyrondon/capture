package payload_test

import (
	"testing"

	"encoding/json"

	"github.com/ifreddyrondon/gocapture/payload"
)

func TestArrayNumberPayloadUnmarshalJSON(t *testing.T) {
	tt := []struct {
		name     string
		payload  []byte
		expected payload.ArrayNumberPayload
		err      error
	}{
		{
			"unmarshal with cap",
			[]byte(`{"cap": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
			nil,
		},
		{
			"unmarshal with captures",
			[]byte(`{"captures": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
			nil,
		},
		{
			"unmarshal with data",
			[]byte(`{"data": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
			nil,
		},
		{
			"unmarshal with payload",
			[]byte(`{"payload": [-78.75, -80.5, -73.75, -70.75, -72]}`),
			[]float64{-78.75, -80.5, -73.75, -70.75, -72},
			nil,
		},
		{
			"unmarshal empty payload",
			[]byte(`{"payload": []}`),
			[]float64{},
			nil,
		},
		{
			"unmarshal empty body",
			[]byte(`{}`),
			nil,
			nil,
		},
		{
			"unmarshal error",
			[]byte(`'`),
			nil,
			payload.ErrorUnmarshalPayload,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := payload.ArrayNumberPayload{}
			err := result.UnmarshalJSON(tc.payload)

			if tc.err != err {
				t.Errorf("Expected get the error '%v'. Got '%v'", tc.err, err)
			}

			// if result expected an error do not check for internal attrs
			if tc.err != nil {
				return
			}

			if len(result) != len(tc.expected) {
				t.Errorf("Expected payload to be '%v'. Got '%v'", len(tc.expected), len(result))
			}

			for i, v := range tc.expected {
				if v != result[i] {
					t.Fatalf("Expected payload at index %v to be '%v'. Got '%v'", i, v, result[i])
				}
			}
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
			result, _ := json.Marshal(tc.payload)

			if tc.expected != string(result) {
				t.Errorf("Expected Marshal result to be '%v'. Got '%v'", tc.expected, string(result))
			}
		})
	}
}
