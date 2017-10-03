package geocoding

import "errors"

var (
	LATError = errors.New("latitude out of boundaries, may range from -90.0 to 90.0")
	LONError = errors.New("longitude out of boundaries, may range from -180.0 to 180.0")
)

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	lat float64
	lng float64
}

// NewPoint returns a valid new Point populated by the passed in latitude (lat) and longitude (lng) values.
//For a valid latitude, longitude pair: -90<=latitude<=+90 and -180<=longitude<=180
func NewPoint(lat float64, lng float64) (*Point, error) {
	if lat < -90 || lat > 90 {
		return nil, LATError
	}

	if lng < -180 || lng > 180 {
		return nil, LONError
	}

	return &Point{lat: lat, lng: lng}, nil
}

// Lat returns the point latitude.
func (p *Point) Lat() float64 {
	return p.lat
}

// Lng returns the point longitude.
func (p *Point) Lng() float64 {
	return p.lng
}
