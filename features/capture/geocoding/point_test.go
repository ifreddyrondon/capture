package geocoding_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/features/capture/geocoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPointSuccess(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		lat, lng float64
	}{
		{"valid with lng upper limit", 75, 180},
		{"valid with lat upper limit", 90, -147.45},
		{"valid with decimals", 77.11112223331, 149.99999999},
		{"valid both upper limits", 90, 180},
		{"valid both lower limits", -90.00000, -180.0000},
		{"valid with just point decimal", 90., 180.},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := geocoding.New(tc.lat, tc.lng)
			require.Nil(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tc.lng, *result.LNG)
			assert.Equal(t, tc.lat, *result.LAT)
		})
	}
}

func TestNewPointFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		lat, lng float64
		err      string
	}{
		{"invalid lat > 90", 95, 280, "latitude out of boundaries, may range from -90.0 to 90.0"},
		{"invalid lat < -95", -95, 280, "latitude out of boundaries, may range from -90.0 to 90.0"},
		{"invalid lng > 180", 75, 280, "longitude out of boundaries, may range from -180.0 to 180.0"},
		{"invalid lng with decimals", 77.11112223331, 249.99999999, "longitude out of boundaries, may range from -180.0 to 180.0"},
		{"invalid lng for 2 decimals points", 90, 180.2, "longitude out of boundaries, may range from -180.0 to 180.0"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			point, err := geocoding.New(tc.lat, tc.lng)
			require.Nil(t, point)
			assert.Contains(t, err.Error(), tc.err)
		})
	}
}

func TestUnmarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payload  []byte
		expected *geocoding.Point
	}{
		{
			"lat and lng",
			[]byte(`{"lat":40.7486, "lng":-73.9864}`),
			getPoint(40.7486, -73.9864),
		},
		{
			"latitude and longitude",
			[]byte(`{"latitude":1, "longitude":1}`),
			getPoint(1, 1),
		},
		{
			"mixed latitude and lng",
			[]byte(`{"latitude":1, "lng":1}`),
			getPoint(1, 1),
		},
		{
			"mixed lat and longitude",
			[]byte(`{"lat":1, "longitude":1}`),
			getPoint(1, 1),
		},
		{
			"with elevation",
			[]byte(`{"lat":1, "longitude":1, "elevation": 1}`),
			func() *geocoding.Point {
				p := getPoint(1, 1)
				elevation := 1.0
				p.Elevation = &elevation
				return p
			}(),
		},
		{
			"with altitude",
			[]byte(`{"lat":1, "longitude":1, "altitude": 1}`),
			func() *geocoding.Point {
				p := getPoint(1, 1)
				elevation := 1.0
				p.Elevation = &elevation
				return p
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := geocoding.Point{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tc.expected.LNG, result.LNG)
			assert.Equal(t, tc.expected.LAT, result.LAT)
			assert.Equal(t, tc.expected.Elevation, result.Elevation)
		})
	}
}

func TestUnmarshalJSONFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		payload []byte
		errs    []string
	}{
		{
			"invalid lat",
			[]byte(`{"lat":100, "lng": 1}`),
			[]string{"latitude out of boundaries, may range from -90.0 to 90.0"},
		},
		{
			"invalid lng",
			[]byte(`{"lat":1, "lng": 190}`),
			[]string{"longitude out of boundaries, may range from -180.0 to 180.0"},
		},
		{
			"invalid json",
			[]byte("`"),
			[]string{"cannot unmarshal json into Point value"},
		},
		{
			"missing lat",
			[]byte(`{"lng":1}`),
			[]string{"latitude must not be blank"},
		},
		{
			"missing lng",
			[]byte(`{"lat":1}`),
			[]string{"longitude must not be blank"},
		},
		{
			"missing latitude",
			[]byte(`{"longitude":1}`),
			[]string{"latitude must not be blank"},
		},
		{
			"missing longitude",
			[]byte(`{"latitude":1}`),
			[]string{"longitude must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := geocoding.Point{}
			err := p.UnmarshalJSON(tc.payload)
			assert.Error(t, err)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

func TestUnmarshalJSONMissingBody(t *testing.T) {
	t.Parallel()

	result := geocoding.Point{}
	err := result.UnmarshalJSON([]byte("{}"))
	require.Nil(t, err)
	require.Nil(t, result.LAT)
	require.Nil(t, result.LNG)
}

func getPoint(lat, lng float64) *geocoding.Point {
	p, _ := geocoding.New(lat, lng)
	return p
}
