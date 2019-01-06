package capture

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

const (
	errLATMissing = "latitude must not be blank"
	errLNGMissing = "longitude must not be blank"
	errLATRange   = "latitude out of boundaries, may range from -90.0 to 90.0"
	errLNGRange   = "longitude out of boundaries, may range from -180.0 to 180.0"
)

const GeolocationValidator validator.StringValidator = "cannot unmarshal json into valid geolocation value"

type GeoLocation struct {
	LAT       *float64 `json:"lat"`
	Latitude  *float64 `json:"latitude"`
	LNG       *float64 `json:"lng"`
	Longitude *float64 `json:"longitude"`
	Elevation *float64 `json:"elevation"`
	Altitude  *float64 `json:"altitude"`
}

func (p *GeoLocation) OK() error {
	e := validate.NewErrors()

	lat := getFloat(p.LAT, p.Latitude)
	validateLatBounds(lat, e)

	lng := getFloat(p.LNG, p.Longitude)
	validateLngBounds(lng, e)

	if lng != nil && lat == nil {
		e.Add("lat", errLATMissing)
	} else if lat != nil && lng == nil {
		e.Add("lng", errLNGMissing)
	}

	if e.HasAny() {
		return e
	}

	p.LAT = lat
	p.LNG = lng
	p.Elevation = getFloat(p.Elevation, p.Altitude)

	return nil
}

func getFloat(data1, data2 *float64) *float64 {
	if data1 == nil {
		return data2
	}
	return data1
}

func validateLatBounds(lat *float64, e *validate.Errors) {
	if lat != nil && !govalidator.InRangeFloat64(*lat, -90, 90) {
		e.Add("lat", errLATRange)
	}
}

func validateLngBounds(lng *float64, e *validate.Errors) {
	if lng != nil && !govalidator.InRangeFloat64(*lng, -180, 180) {
		e.Add("lat", errLNGRange)
	}
}
