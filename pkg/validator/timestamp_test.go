package validator_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/validator"
)

func s2t(date string) *time.Time {
	v, _ := time.Parse(time.RFC3339, date)
	return &v
}

func TestValidateTimestampOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name              string
		body              string
		expectedTimestamp *time.Time
	}{
		{
			"decode timestamp without data",
			`{}`,
			nil,
		},
		{
			"decode timestamp with date (RFC3339)",
			`{"date":"1989-12-26T06:01:00.00Z"}`,
			s2t("1989-12-26T06:01:00.00Z"),
		},
		{
			"decode timestamp with timestamp (RFC3339)",
			`{"timestamp":"1989-12-26T06:01:00.00Z"}`,
			s2t("1989-12-26T06:01:00.00Z"),
		},
		{
			"decode RFC1123 date",
			`{"date": "Tue, 26 Dec 1989 06:01:00 UTC"}`,
			s2t("1989-12-26T06:01:00.00Z"),
		},
		{
			"decode Unix timestamp as string date",
			`{"date": "630655260"}`,
			s2t("1989-12-26T03:01:00-03:00"),
		},
		{
			"decode Unix timestamp as integer date",
			`{"date": 630655260}`,
			s2t("1989-12-26T03:01:00-03:00"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var d validator.Timestamp
			err := validator.TimestampValidator.Decode(r, &d)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedTimestamp, d.Time)
		})
	}
}

func TestDecodeTimestampFails(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			"decode timestamp when invalid json",
			".",
			"cannot unmarshal json into valid time value",
		},
		{
			"decode timestamp when invalid date",
			`{"date": "asd"}`,
			"Could not find date format for asd",
		},
		{
			"decode timestamp when invalid timestamp",
			`{"timestamp": "asd"}`,
			"Could not find date format for asd",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var d validator.Timestamp
			err := validator.TimestampValidator.Decode(r, &d)
			assert.EqualError(t, err, tc.err)
		})
	}
}
