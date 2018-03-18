package geocoding

import (
	"log"

	jlexer "github.com/mailru/easyjson/jlexer"
)

type pointJSON struct {
	Lat       float64 `json:"lat"`
	Latitude  float64 `json:"latitude"`
	Lng       float64 `json:"lng"`
	Longitude float64 `json:"longitude"`
}

func (p *pointJSON) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3844eb60DecodeGithubComIfreddyrondonGocaptureGeocoding(&r, p)
	if err := r.Error(); err != nil {
		log.Print(err)
		return ErrorUnmarshalPoint
	}

	return nil
}
