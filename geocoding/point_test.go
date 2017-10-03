package geocoding_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/geocoding"
)

func TestNewPoint(t *testing.T) {
	tt := []struct {
		name     string
		lat, lng float64
		err      error
	}{
		{"valid with lng upper limit", 75, 180, nil},
		{"valid with lat upper limit", 90, -147.45, nil},
		{"valid with decimals", 77.11112223331, 149.99999999, nil},
		{"valid both upper limits", 90, 180, nil},
		{"valid both lower limits", -90.00000, -180.0000, nil},
		{"valid with just point decimal", 90., 180., nil},
		{"invalid lat", 75, 280, geocoding.LATError},
		{"invalid lng", 75, 280, geocoding.LONError},
		{"invalid lng with decimals", 77.11112223331, 249.99999999, geocoding.LONError},
		{"invalid lng for 2 decimals points", 90, 180.2, geocoding.LONError},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := geocoding.NewPoint(tc.lat, tc.lng)

			if tc.err == nil {
				if p == nil {
					t.Errorf("Expected point not to nil. Got '%v'", p)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error to be %v. Got '%v'", tc.err, p)
				}
			}
		})
	}
}
