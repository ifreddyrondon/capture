package geocoding

import (
	"errors"
	"log"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

var (
	ErrorLATRange       = errors.New("latitude out of boundaries, may range from -90.0 to 90.0")
	ErrorLONRange       = errors.New("longitude out of boundaries, may range from -180.0 to 180.0")
	ErrorLATMissing     = errors.New("missing latitude")
	ErrorLNGMissing     = errors.New("missing longitude")
	ErrorUnmarshalPoint = errors.New("cannot unmarshal json into Point value")
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371
)

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// NewPoint returns a valid new Point populated by the passed in latitude (lat) and longitude (lng) values.
// For a valid latitude, longitude pair: -90<=latitude<=+90 and -180<=longitude<=180
func NewPoint(lat float64, lng float64) (*Point, error) {
	if lat < -90 || lat > 90 {
		return nil, ErrorLATRange
	}

	if lng < -180 || lng > 180 {
		return nil, ErrorLONRange
	}

	return &Point{Lat: lat, Lng: lng}, nil
}

// MarshalJSON decode current Point to JSON.
// It supports json.Marshaler interface.
func (po Point) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3844eb60EncodeGithubComIfreddyrondonGocaptureGeocoding1(&w, po)
	return w.Buffer.BuildBytes(), w.Error
}

type pointJSON struct {
	Lat       float64 `json:"lat"`
	Latitude  float64 `json:"latitude"`
	Lng       float64 `json:"lng"`
	Longitude float64 `json:"longitude"`
}

// UnmarshalJSON decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (po *Point) UnmarshalJSON(data []byte) error {
	model := new(pointJSON)
	r := jlexer.Lexer{Data: data}
	easyjson3844eb60DecodeGithubComIfreddyrondonGocaptureGeocoding(&r, model)
	if err := r.Error(); err != nil {
		log.Print(err)
		return ErrorUnmarshalPoint
	}

	lat, lng := getLat(model), getLng(model)
	if lat == 0 {
		return ErrorLATMissing
	}
	if lng == 0 {
		return ErrorLNGMissing
	}

	point, err := NewPoint(lat, lng)
	if err != nil {
		log.Print(err)
		return err
	}

	*po = *point

	return nil
}

func getLat(model *pointJSON) float64 {
	var lat float64
	if model.Lat != 0 {
		lat = model.Lat
	} else if model.Latitude != 0 {
		lat = model.Latitude
	}
	return lat
}

func getLng(model *pointJSON) float64 {
	var lng float64
	if model.Lng != 0 {
		lng = model.Lng
	} else if model.Longitude != 0 {
		lng = model.Longitude
	}
	return lng
}
