package geocoding

import (
	"errors"
	"log"

	jwriter "github.com/mailru/easyjson/jwriter"
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
	LAT *float64 `json:"lat"`
	LNG *float64 `json:"lng"`
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

	return &Point{LAT: &lat, LNG: &lng}, nil
}

// MarshalJSON decode current Point to JSON.
// It supports json.Marshaler interface.
func (p Point) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3844eb60EncodeGithubComIfreddyrondonGocaptureGeocoding(&w, p)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (p *Point) UnmarshalJSON(data []byte) error {
	var model pointJSON
	if err := model.unmarshalJSON(data); err != nil {
		return err
	}

	lat := getLAT(&model)
	lng := getLNG(&model)
	if lat == nil && lng == nil {
		return nil
	}
	if lng != nil && lat == nil {
		return ErrorLATMissing
	}
	if lat != nil && lng == nil {
		return ErrorLNGMissing
	}

	point, err := New(*lat, *lng)
	if err != nil {
		log.Print(err)
		return err
	}

	*p = *point

	return nil
}

func getLAT(model *pointJSON) *float64 {
	if model.LAT == nil {
		return model.Latitude
	}
	return model.LAT
}

func getLNG(model *pointJSON) *float64 {
	if model.LNG == nil {
		return model.Longitude
	}
	return model.LNG
}
