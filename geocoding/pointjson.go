package geocoding

import (
	jlexer "github.com/mailru/easyjson/jlexer"
)

type pointJSON struct {
	LAT       *float64 `json:"lat"`
	Latitude  *float64 `json:"latitude"`
	LNG       *float64 `json:"lng"`
	Longitude *float64 `json:"longitude"`
}

func (p *pointJSON) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA2deb046DecodeGithubComIfreddyrondonGocaptureGeocoding(&r, p)
	if err := r.Error(); err != nil {
		return ErrorUnmarshalPoint
	}
	return nil
}
