package numberlist

import (
	"errors"

	"log"

	"github.com/mailru/easyjson/jlexer"
)

var (
	ErrorUnmarshalPayload = errors.New("cannot unmarshal json into Payload valid value")
)

// Payload represent an association of float numbers
type Payload []float64

// New returns a new pointer to a ArrayNumberPayload composed of the passed float64
func New(data ...float64) *Payload {
	p := new(Payload)
	*p = data
	return p
}

type jsonPayload struct {
	Cap      []float64 `json:"cap"`
	Captures []float64 `json:"captures"`
	Data     []float64 `json:"data"`
	Payload  []float64 `json:"payload"`
}

// UnmarshalJSON supports json.Unmarshaler interface
func (p *Payload) UnmarshalJSON(data []byte) error {
	model := new(jsonPayload)
	if err := model.unmarshalJSON(data); err != nil {
		log.Print(err)
		return ErrorUnmarshalPayload
	}
	*p = model.getPayload()
	return nil
}

func (v *jsonPayload) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComIfreddyrondonGocapturePayload(&r, v)
	return r.Error()
}

func (v *jsonPayload) getPayload() []float64 {
	if v.Cap != nil {
		return v.Cap
	} else if v.Captures != nil {
		return v.Captures
	} else if v.Data != nil {
		return v.Data
	}
	return v.Payload
}
