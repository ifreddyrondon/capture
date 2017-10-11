package geocoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
)

var (
	RangeLATError       = errors.New("latitude out of boundaries, may range from -90.0 to 90.0")
	RangeLONError       = errors.New("longitude out of boundaries, may range from -180.0 to 180.0")
	MissingLATError     = errors.New("missing latitude")
	MissingLNGError     = errors.New("missing longitude")
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
		return nil, RangeLATError
	}

	if lng < -180 || lng > 180 {
		return nil, RangeLONError
	}

	return &Point{Lat: lat, Lng: lng}, nil
}

// UnmarshalJSON decode a JSON body into a Point value
// Throws an error if the payload cannot be interpreted as a JSON body
func UnmarshalJSON(body []byte) (*Point, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	var values map[string]float64
	if err := decoder.Decode(&values); err != nil {
		log.Print(err)
		return nil, PointUnmarshalError
	}

	var lat, lng float64
	var isOk bool
	if lat, isOk = getLatitude(&values); !isOk {
		return nil, MissingLATError
	}
	if lng, isOk = getLongitude(&values); !isOk {
		return nil, MissingLNGError
	}

	p, err := NewPoint(lat, lng)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return p, nil
}

func getLatitude(values *map[string]float64) (float64, bool) {
	var val float64
	var isOk bool
	if val, isOk = (*values)["lat"]; isOk {
		isOk = true
	} else if val, isOk = (*values)["latitude"]; isOk {
		isOk = true
	}

	return val, isOk
}

func getLongitude(values *map[string]float64) (float64, bool) {
	var val float64
	var isOk bool
	if val, isOk = (*values)["lng"]; isOk {
		isOk = true
	} else if val, isOk = (*values)["longitude"]; isOk {
		isOk = true
	}

	return val, isOk
}
