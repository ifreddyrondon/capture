package timestamp_test

import (
	"testing"

	"time"

	"github.com/ifreddyrondon/gocapture/timestamp"
)

func TestNewDate(t *testing.T) {
	date := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	result := timestamp.NewTimestamp(date)
	expetedResult := "1989-12-26 06:01:00 +0000 UTC"

	if result == nil {
		t.Errorf("Expected capture not to nil. Got '%v'", result)
	}

	if result.Timestamp.String() != expetedResult {
		t.Errorf("Expected date to be '%v'. Got '%v'", expetedResult, result.Timestamp.String())
	}
}

func TestUnmarshalJSON(t *testing.T) {
	expectResult := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
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
			result := &timestamp.Timestamp{}
			result.UnmarshalJSON(tc.payload)

			if !result.Timestamp.Equal(expectResult) {
				t.Errorf("Expected CaptureDate to be '%v'. Got '%v'", expectResult, result.Timestamp)
			}
		})
	}
}

func TestUnmarshalJSONWhenFails(t *testing.T) {
	expectResult := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	fakeTime := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	mockClock := timestamp.NewMockClock(fakeTime)

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
			result := &timestamp.Timestamp{}
			timestamp.SetClockInstance(result, mockClock)
			result.UnmarshalJSON(tc.payload)

			if !result.Timestamp.Equal(expectResult) {
				t.Errorf("Expected CaptureDate to be '%v'. Got '%v'", expectResult, result.Timestamp)
			}
		})
	}
}

func TestDateMarshalJSON(t *testing.T) {
	parsedDate, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := `"1989-12-26T06:01:00Z"`
	result, _ := timestamp.NewTimestamp(parsedDate).MarshalJSON()
	if string(result) != expected {
		t.Errorf("Expected Marshall data to be '%v'. Got '%v'", expected, string(result))
	}
}
