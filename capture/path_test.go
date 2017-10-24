package capture_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/capture"
)

func TestPathAddCapture(t *testing.T) {
	tt := []struct {
		name     string
		captures []*capture.Capture
	}{
		{"empty path", []*capture.Capture{}},
		{"empty path", []*capture.Capture{
			getCapture(1, 1, "1989-12-26T06:01:00.00Z", ""),
			getCapture(1, 2, "1989-12-26T06:01:00.00Z", ""),
		}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := new(capture.Path)
			p.AddCapture(tc.captures...)

			if p == nil {
				t.Errorf("Expected path not to nil. Got '%v'", p)
			}

			if len(p.Captures) != len(tc.captures) {
				t.Errorf("Expected path len to be '%v'. Got '%v'", len(tc.captures), len(p.Captures))
			}
		})
	}
}
