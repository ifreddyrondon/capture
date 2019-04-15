package validator_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

func TestValidatePayloadOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     string
		expected validator.Payload
	}{
		{
			name: "decode payload with data",
			body: `{"data": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}]}`,
			expected: validator.Payload{Payload: []domain.Metric{
				{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			}},
		},
		{
			name: "decode payload with payload",
			body: `{"payload": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}]}`,
			expected: validator.Payload{Payload: []domain.Metric{
				{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
			}},
		},
		{
			name: "decode payload with payload and two metrics as raw data",
			body: `{"payload": [{"name": "temp", "value": 10}, {"name": "power", "value": 30}]}`,
			expected: validator.Payload{Payload: []domain.Metric{
				{Name: "temp", Value: 10.0},
				{Name: "power", Value: 30.0},
			}},
		},
		{
			name: "decode payload with payload and two metrics as raw array",
			body: `{"payload": [{"name": "power", "value": [-78.75, -80.5, -73.75, -70.75, -72]}, {"name": "frequencies", "value": [100, 200, 300, 400, 500]}]}`,
			expected: validator.Payload{Payload: []domain.Metric{
				{Name: "power", Value: []interface{}{-78.75, -80.5, -73.75, -70.75, -72.0}},
				{Name: "frequencies", Value: []interface{}{100.0, 200.0, 300.0, 400.0, 500.0}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p validator.Payload
			err := validator.PayloadValidator.Decode(r, &p)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Payload, p.Payload)
		})
	}
}

func TestValidatePayloadFails(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			"unmarshal error",
			`'`,
			"cannot unmarshal json into valid payload value",
		},
		{
			"unmarshal nil payload",
			`{"payload": null`,
			"cannot unmarshal json into valid payload value",
		},
		{
			"unmarshal empty payload",
			`{"payload": []`,
			"cannot unmarshal json into valid payload value",
		},
		{
			"unmarshal payload with nulls",
			`{"payload": [null]`,
			"cannot unmarshal json into valid payload value",
		},
		{
			"unmarshal empty body",
			`{}`,
			"payload value must not be blank",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p validator.Payload
			err := validator.PayloadValidator.Decode(r, &p)
			assert.EqualError(t, err, tc.err)
		})
	}
}
