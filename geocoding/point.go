package geocoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
)

var (
	LATError            = errors.New("latitude out of boundaries, may range from -90.0 to 90.0")
	LONError            = errors.New("longitude out of boundaries, may range from -180.0 to 180.0")
	PointUnmarshalError = errors.New("cannot unmarshal json into Point value")
)

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
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

	return &Point{Lat: lat, Lng: lng}, nil
}

// UnmarshalJSON decode a JSON body into a Point value
// Throws an error if the body of the point cannot be interpreted by the JSON body
func UnmarshalJSON(body []byte) (*Point, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	var values map[string]float64
	if err := decoder.Decode(&values); err != nil {
		log.Print(err)
		return nil, PointUnmarshalError
	}

	p, err := NewPoint(values["lat"], values["lng"])
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return p, nil
}
