package geocoding

import (
	"errors"

	"github.com/asaskevich/govalidator"

	"github.com/markbates/validate"
)

// ErrUnmarshalPoint expected error when fails to unmarshal a point
var ErrUnmarshalPoint = errors.New("cannot unmarshal json into Point value")

const (
	// errLATRange expected error when latitude is out of boundaries
	errLATRange = "latitude out of boundaries, may range from -90.0 to 90.0"
	// errLNGRange expected error when longitude is out of boundaries
	errLNGRange = "longitude out of boundaries, may range from -180.0 to 180.0"
)

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	LAT       *float64 `json:"lat"`
	LNG       *float64 `json:"lng"`
	Elevation *float64 `json:"elevation"`
}

// New returns a valid new Point populated by the passed in latitude (lat) and longitude (lng) values.
// For a valid latitude, longitude pair: -90<=latitude<=+90 and -180<=longitude<=180
func New(lat float64, lng float64) (*Point, error) {
	p := &Point{LAT: &lat, LNG: &lng}
	if err := validate.Validate(p); err.HasAny() {
		return nil, err
	}
	return p, nil
}

// IsValid validates if a point is valid.
// It must be -90<=latitude<=+90 and -180<=longitude<=180
func (p *Point) IsValid(errors *validate.Errors) {
	if p.LAT != nil && !govalidator.InRangeFloat64(*p.LAT, -90, 90) {
		errors.Add("lat", errLATRange)
	}
	if p.LNG != nil && !govalidator.InRangeFloat64(*p.LNG, -180, 180) {
		errors.Add("lat", errLNGRange)
	}
}

// UnmarshalJSON decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (p *Point) UnmarshalJSON(data []byte) error {
	var model pointJSON
	if err := model.unmarshalJSON(data); err != nil {
		return err
	}

	if err := validate.Validate(&model); err.HasAny() {
		return err
	}
	point := model.getPoint()
	if err := validate.Validate(point); err.HasAny() {
		return err
	}

	*p = *point

	return nil
}
