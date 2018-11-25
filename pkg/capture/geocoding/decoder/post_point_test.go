package decoder_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/capture/geocoding"
	"github.com/ifreddyrondon/capture/pkg/capture/geocoding/decoder"
	"github.com/stretchr/testify/assert"
)

func f2P(v float64) *float64 {
	return &v
}

func TestDecodePostPointOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     string
		expected decoder.PostPoint
	}{
		{
			name:     "decode point without data",
			body:     `{}`,
			expected: decoder.PostPoint{},
		},
		{
			name:     "decode point with lat, lng and without elevation",
			body:     `{"lat":75,"lng":180}`,
			expected: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with lat, longitude and without elevation",
			body:     `{"lat":75,"longitude":180}`,
			expected: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with latitude, lng and without elevation",
			body:     `{"latitude":75,"lng":180}`,
			expected: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with latitude, longitude and without elevation",
			body:     `{"latitude":75,"longitude":180}`,
			expected: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with lat, lng and elevation",
			body:     `{"lat":75,"lng":180,"elevation":1}`,
			expected: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180), Elevation: f2P(1)},
		},
		{
			name:     "decode point with lat, lng and altitude",
			body:     `{"lat":75,"lng":180,"altitude":1}`,
			expected: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180), Elevation: f2P(1)},
		},
		{
			name:     "decode point with lat, lng for upper limits",
			body:     `{"lat":90,"lng":-147.45}`,
			expected: decoder.PostPoint{LAT: f2P(90), LNG: f2P(-147.45)},
		},
		{
			name:     "decode point with lat, lng for lower limits",
			body:     `{"lat":-90.00000,"lng":-180.0000}`,
			expected: decoder.PostPoint{LAT: f2P(-90.00000), LNG: f2P(-180.0000)},
		},
		{
			name:     "decode point with lat, lng with decimals",
			body:     `{"lat":77.11112223331,"lng":149.99999999}`,
			expected: decoder.PostPoint{LAT: f2P(77.11112223331), LNG: f2P(149.99999999)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p decoder.PostPoint
			err := decoder.Decode(r, &p)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.LAT, p.LAT)
			assert.Equal(t, tc.expected.LNG, p.LNG)
			assert.Equal(t, tc.expected.Elevation, p.Elevation)

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
			"decode point when missing lat",
			`{"lng":180}`,
			"latitude must not be blank",
		},
		{
			"decode point when missing lng",
			`{"lat":1}`,
			"longitude must not be blank",
		},
		{
			"decode point when missing latitude",
			`{"longitude":1}`,
			"latitude must not be blank",
		},
		{
			"decode point when missing longitude",
			`{"latitude":1}`,
			"longitude must not be blank",
		},
		{
			"decode point when invalid lat",
			`{"lat":100, "lng": 1}`,
			"latitude out of boundaries, may range from -90.0 to 90.0",
		},
		{
			"decode point when invalid lng",
			`{"lat":1, "lng": 190}`,
			"longitude out of boundaries, may range from -180.0 to 180.0",
		},
		{
			"decode point when invalid json",
			".",
			"cannot unmarshal json into point value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p decoder.PostPoint
			err := decoder.Decode(r, &p)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestPointFromPostPointOK(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name      string
		postPoint decoder.PostPoint
		expected  geocoding.Point
	}{
		{
			name:      "get point from postPoint with lat and lng",
			postPoint: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180)},
			expected:  geocoding.Point{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:      "get point from postPoint with lat, lng and elevation",
			postPoint: decoder.PostPoint{LAT: f2P(75), LNG: f2P(180), Elevation: f2P(1)},
			expected:  geocoding.Point{LAT: f2P(75), LNG: f2P(180), Elevation: f2P(1)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.postPoint.GetPoint()
			assert.Equal(t, tc.expected.LAT, p.LAT)
			assert.Equal(t, tc.expected.LNG, p.LNG)
			assert.Equal(t, tc.expected.Elevation, p.Elevation)
		})
	}
}
