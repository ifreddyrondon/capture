package geocoding_test

import (
	"testing"

	"fmt"

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
		{"invalid lat > 90", 95, 280, geocoding.RangeLATError},
		{"invalid lat < -95", -95, 280, geocoding.RangeLATError},
		{"invalid lng > 180", 75, 280, geocoding.RangeLONError},
		{"invalid lng with decimals", 77.11112223331, 249.99999999, geocoding.RangeLONError},
		{"invalid lng for 2 decimals points", 90, 180.2, geocoding.RangeLONError},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := geocoding.NewPoint(tc.lat, tc.lng)

			if tc.err == nil {
				if p == nil {
					t.Errorf("Expected point not to nil. Got '%v'", p)
				}

				if p.Lat != tc.lat {
					t.Errorf("Expected result lat point to be '%v'. Got '%v'", tc.lat, p.Lat)
				}

				if p.Lng != tc.lng {
					t.Errorf("Expected result lng point to be '%v'. Got '%v'", tc.lng, p.Lng)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error to be %v. Got '%v'", tc.err, p)
				}
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tt := []struct {
		name        string
		payload     []byte
		resultPoint *geocoding.Point
		resultError error
	}{
		{
			"valid lat and lng",
			pointToBytes(40.7486, -73.9864),
			getPoint(40.7486, -73.9864),
			nil,
		},
		{
			"valid with latitude and longitude",
			[]byte(`{"latitude":1, "longitude":1}`),
			getPoint(1, 1),
			nil,
		},
		{
			"valid mixed latitude and lng",
			[]byte(`{"latitude":1, "lng":1}`),
			getPoint(1, 1),
			nil,
		},
		{
			"valid mixed lat and longitude",
			[]byte(`{"lat":1, "longitude":1}`),
			getPoint(1, 1),
			nil,
		},
		{"invalid empty json", []byte("{}"), nil, geocoding.MissingLATError},
		{"invalid lat", pointToBytes(100, 1), nil, geocoding.RangeLATError},
		{"invalid lng", pointToBytes(1, 190), nil, geocoding.RangeLONError},
		{"invalid json", []byte("`"), nil, geocoding.PointUnmarshalError},
		{"missing lat", []byte(`{"lng":1}`), nil, geocoding.MissingLATError},
		{"missing lng", []byte(`{"lat":1}`), nil, geocoding.MissingLNGError},
		{"missing latitude", []byte(`{"longitude":1}`), nil, geocoding.MissingLATError},
		{"missing longitude", []byte(`{"latitude":1}`), nil, geocoding.MissingLNGError},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resultPoint, resultError := geocoding.UnmarshalJSON(tc.payload)

			if resultError != tc.resultError {
				t.Errorf("Expected get the error '%v'. Got '%v'", tc.resultError, resultError)
			}

			// if result expected an error do not check for internal attrs
			if tc.resultError != nil {
				return
			}

			if resultPoint.Lat != tc.resultPoint.Lat {
				t.Errorf("Expected result lat point to be '%v'. Got '%v'", tc.resultPoint.Lat, resultPoint.Lat)
			}

			if resultPoint.Lng != tc.resultPoint.Lng {
				t.Errorf("Expected result lng point to be '%v'. Got '%v'", tc.resultPoint.Lng, resultPoint.Lng)
			}
		})
	}
}

func pointToBytes(lat, lng float64) []byte {
	res := fmt.Sprintf(`{"lat":%v, "lng":%v}`, lat, lng)
	return []byte(res)
}

func getPoint(lat, lng float64) *geocoding.Point {
	p, _ := geocoding.NewPoint(lat, lng)
	return p
}
