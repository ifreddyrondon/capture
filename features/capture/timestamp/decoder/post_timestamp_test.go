package decoder

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/stretchr/testify/assert"
)

var defaultTimeNow = time.Date(2017, time.Month(12), 26, 6, 1, 0, 0, time.UTC)

func s2t(date string) *time.Time {
	v, _ := time.Parse(time.RFC3339, date)
	return &v
}

func TestDecodePostTimestampOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name                  string
		body                  string
		expectedPostTimestamp *time.Time
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

			var d PostTimestamp
			err := Decode(r, &d)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedPostTimestamp, d.postTimestamp)
		})
	}
}

func TestDecodePostPointError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			"decode timestamp when invalid json",
			".",
			"cannot unmarshal json into timestamp value",
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

			var d PostTimestamp
			err := Decode(r, &d)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestPointFromPostTimestampOK(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name      string
		postPoint PostTimestamp
		expected  time.Time
	}{
		{
			"get time from PostTimestamp without data",
			PostTimestamp{},
			*s2t("2017-12-26T06:01:00.00Z"),
		},
		{
			"get time from PostTimestamp with data",
			PostTimestamp{postTimestamp: s2t("1989-12-26T06:01:00.00Z")},
			*s2t("1989-12-26T06:01:00.00Z"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.postPoint.clock = features.NewMockClock(defaultTimeNow)
			date := tc.postPoint.GetTimestamp()
			assert.Equal(t, tc.expected, date)
		})
	}
}
