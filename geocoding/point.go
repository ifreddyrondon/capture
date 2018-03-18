package geocoding

import (
	"errors"
	"log"

	"github.com/markbates/going/defaults"

	"github.com/mailru/easyjson/jwriter"
)

var (
	// ErrorLATRange expected error when latitude is out of boundaries
	ErrorLATRange = errors.New("latitude out of boundaries, may range from -90.0 to 90.0")
	// ErrorLONRange expected error when longitude is out of boundaries
	ErrorLONRange = errors.New("longitude out of boundaries, may range from -180.0 to 180.0")
	// ErrorLATMissing expected error when latitude is missing
	ErrorLATMissing = errors.New("missing latitude")
	// ErrorLNGMissing expected error when longitude is missing
	ErrorLNGMissing = errors.New("missing longitude")
	// ErrorUnmarshalPoint expected error when fails to unmarshal a point
	ErrorUnmarshalPoint = errors.New("cannot unmarshal json into Point value")
)

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// New returns a valid new Point populated by the passed in latitude (lat) and longitude (lng) values.
// For a valid latitude, longitude pair: -90<=latitude<=+90 and -180<=longitude<=180
func New(lat float64, lng float64) (*Point, error) {
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

// UnmarshalJSON decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (po *Point) UnmarshalJSON(data []byte) error {
	var model pointJSON
	if err := model.unmarshalJSON(data); err != nil {
		return err
	}

	lat := defaults.Float64(model.Lat, model.Latitude)
	if lat == 0 {
		return ErrorLATMissing
	}
	lng := defaults.Float64(model.Lng, model.Longitude)
	if lng == 0 {
		return ErrorLNGMissing
	}

	point, err := New(lat, lng)
	if err != nil {
		log.Print(err)
		return err
	}

	*po = *point

	return nil
}
