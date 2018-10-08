package geocoding

import (
	jlexer "github.com/mailru/easyjson/jlexer"
	"github.com/markbates/validate"
)

const (
	errLATMissing = "latitude must not be blank"
	errLNGMissing = "longitude must not be blank"
)

type pointJSON struct {
	LAT       *float64 `json:"lat"`
	Latitude  *float64 `json:"latitude"`
	LNG       *float64 `json:"lng"`
	Longitude *float64 `json:"longitude"`
	Elevation *float64 `json:"elevation"`
	Altitude  *float64 `json:"altitude"`
}

func (p *pointJSON) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA2deb046DecodeGithubComIfreddyrondonCaptureGeocoding(&r, p)
	if err := r.Error(); err != nil {
		return ErrUnmarshalPoint
	}
	return nil
}

func (p *pointJSON) getLat() *float64 {
	return getFloat(p.LAT, p.Latitude)
}

func (p *pointJSON) getLng() *float64 {
	return getFloat(p.LNG, p.Longitude)
}

func (p *pointJSON) getElevation() *float64 {
	return getFloat(p.Elevation, p.Altitude)
}

func (p *pointJSON) getPoint() *Point {
	return &Point{
		LAT:       p.getLat(),
		LNG:       p.getLng(),
		Elevation: p.getElevation(),
	}
}

func (p *pointJSON) IsValid(errors *validate.Errors) {
	lat := p.getLat()
	lng := p.getLng()

	if lng != nil && lat == nil {
		errors.Add("lat", errLATMissing)
	} else if lat != nil && lng == nil {
		errors.Add("lng", errLNGMissing)
	}
}

func getFloat(data1, data2 *float64) *float64 {
	if data1 == nil {
		return data2
	}
	return data1
}
