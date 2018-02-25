package timestamp_test

import (
	"testing"

	"time"

	"github.com/ifreddyrondon/gocapture/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestNewDate(t *testing.T) {
	date := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	ts := timestamp.NewTimestamp(date)
	expected := "1989-12-26 06:01:00 +0000 UTC"
	assert.Equal(t, expected, ts.Timestamp.String())
}

func TestUnmarshalJSON(t *testing.T) {
	expected := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	tt := []struct {
		name    string
		payload []byte
	}{
		{
			"valid RFC3339 date with date key",
			[]byte(`{"date": "1989-12-26T06:01:00.00Z"}`),
		},
		{
			"valid RFC3339 date with timestamp key",
			[]byte(`{"timestamp": "1989-12-26T06:01:00.00Z"}`),
		},
		{
			"valid RFC1123 date with date key",
			[]byte(`{"date": "Tue, 26 Dec 1989 06:01:00 UTC"}`),
		},
		{
			"valid RFC1123 date with timestamp key",
			[]byte(`{"timestamp": "Tue, 26 Dec 1989 06:01:00 UTC"}`),
		},
		{
			"valid Unix timestamp as string date with date key",
			[]byte(`{"date": "630655260"}`),
		},
		{
			"valid Unix timestamp as string date with timestamp key",
			[]byte(`{"timestamp": "630655260"}`),
		},
		{
			"valid Unix timestamp as integer date with date key",
			[]byte(`{"date": 630655260}`),
		},
		{
			"valid Unix timestamp as integer date with timestamp key",
			[]byte(`{"timestamp": 630655260}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := &timestamp.Timestamp{}
			ts.UnmarshalJSON(tc.payload)
			assert.Equal(t, expected, ts.Timestamp)
		})
	}
}

func TestUnmarshalJSONWhenFails(t *testing.T) {
	expected := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	mockClock := timestamp.NewMockClock(expected)

	tt := []struct {
		name    string
		payload []byte
	}{
		{"empty json", []byte("{}")},
		{"invalid json", []byte("`")},
		{"missing date or timestamp", []byte(`{"foo": "630655260"}`)},
		{"invalid date", []byte(`{"date": "asd"}`)},
		{"invalid timestamp", []byte(`{"timestamp": "asd"}`)},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := &timestamp.Timestamp{}
			timestamp.SetClockInstance(ts, mockClock)
			ts.UnmarshalJSON(tc.payload)
			assert.Equal(t, expected, ts.Timestamp)
		})
	}
}

func TestDateMarshalJSON(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")
	expected := `"1989-12-26T06:01:00Z"`
	result, _ := timestamp.NewTimestamp(date).MarshalJSON()
	assert.Equal(t, expected, string(result))
}