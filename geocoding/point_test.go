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

func TestPoint_Lat(t *testing.T) {
	p, _ := geocoding.NewPoint(40.5, 120.5)

	lat := p.Lat()

	if lat != 40.5 {
		t.Errorf("Expected Lat to be '%v'. Got '%v'", 40.5, lat)
	}
}

// Tests that calling GetLng() after creating a new point returns the expected lng value.
func TestPoint_Lng(t *testing.T) {
	p, _ := geocoding.NewPoint(40.5, 120.5)

	lng := p.Lng()

	if lng != 120.5 {
		t.Errorf("Expected Lat to be '%v'. Got '%v'", 120.5, lng)
	}
}
