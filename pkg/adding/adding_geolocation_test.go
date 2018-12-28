package adding_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/stretchr/testify/assert"
)

func f2P(v float64) *float64 {
	return &v
}

func TestValidateGeolocationOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     string
		expected adding.GeoLocation
	}{
		{
			name:     "decode point without data",
			body:     `{}`,
			expected: adding.GeoLocation{},
		},
		{
			name:     "decode point with lat, lng and without elevation",
			body:     `{"lat":75,"lng":180}`,
			expected: adding.GeoLocation{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with lat, longitude and without elevation",
			body:     `{"lat":75,"longitude":180}`,
			expected: adding.GeoLocation{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with latitude, lng and without elevation",
			body:     `{"latitude":75,"lng":180}`,
			expected: adding.GeoLocation{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with latitude, longitude and without elevation",
			body:     `{"latitude":75,"longitude":180}`,
			expected: adding.GeoLocation{LAT: f2P(75), LNG: f2P(180)},
		},
		{
			name:     "decode point with lat, lng and elevation",
			body:     `{"lat":75,"lng":180,"elevation":1}`,
			expected: adding.GeoLocation{LAT: f2P(75), LNG: f2P(180), Elevation: f2P(1)},
		},
		{
			name:     "decode point with lat, lng and altitude",
			body:     `{"lat":75,"lng":180,"altitude":1}`,
			expected: adding.GeoLocation{LAT: f2P(75), LNG: f2P(180), Elevation: f2P(1)},
		},
		{
			name:     "decode point with lat, lng for upper limits",
			body:     `{"lat":90,"lng":-147.45}`,
			expected: adding.GeoLocation{LAT: f2P(90), LNG: f2P(-147.45)},
		},
		{
			name:     "decode point with lat, lng for lower limits",
			body:     `{"lat":-90.00000,"lng":-180.0000}`,
			expected: adding.GeoLocation{LAT: f2P(-90.00000), LNG: f2P(-180.0000)},
		},
		{
			name:     "decode point with lat, lng with decimals",
			body:     `{"lat":77.11112223331,"lng":149.99999999}`,
			expected: adding.GeoLocation{LAT: f2P(77.11112223331), LNG: f2P(149.99999999)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p adding.GeoLocation
			err := adding.GeolocationValidator.Decode(r, &p)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.LAT, p.LAT)
			assert.Equal(t, tc.expected.LNG, p.LNG)
			assert.Equal(t, tc.expected.Elevation, p.Elevation)

		})
	}
}

func TestValidateGeolocationFails(t *testing.T) {
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
			"cannot unmarshal json into valid geolocation value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p adding.GeoLocation
			err := adding.GeolocationValidator.Decode(r, &p)
			assert.EqualError(t, err, tc.err)
		})
	}
}
