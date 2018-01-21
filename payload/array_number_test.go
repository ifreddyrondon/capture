package payload_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/payload"
)

func TestArrayNumberPayloadValues(t *testing.T) {
	expected := []float64{-78.75, -80.5, -73.75, -70.75, -72}

	values := payload.ArrayNumberPayload{-78.75, -80.5, -73.75, -70.75, -72}.Values()
	for i, v := range values.(payload.ArrayNumberPayload) {
		if v != expected[i] {
			t.Fatalf("Expected payload at index %v to be '%v'. Got '%v'", i, expected[i], v)
		}
	}
}

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
