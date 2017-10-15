package captures_test

import (
	"testing"
	"time"

	"github.com/ifreddyrondon/gocapture/captures"
	"github.com/ifreddyrondon/gocapture/geocoding"
)

func TestNewCapture(t *testing.T) {
	point, _ := geocoding.NewPoint(1, 2)
	timestamp := captures.NewDate(time.Now())
	var payload interface{}

	result := captures.NewCapture(point, timestamp, payload)

	if result == nil {
		t.Errorf("Expected capture not to nil. Got '%v'", result)
	}
}
