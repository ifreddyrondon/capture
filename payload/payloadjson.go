package payload

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/markbates/validate"
)

// errMissingPayload expected error when payload is missing
const errMissingPayload = "payload value must not be blank"

type jsonPayload struct {
	Cap      []*Metric `json:"cap"`
	Captures []*Metric `json:"captures"`
	Data     []*Metric `json:"data"`
	Payload  []*Metric `json:"payload"`
}

// UnmarshalJSON supports json.Unmarshaler interface
func (p *jsonPayload) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ad23cceDecodeGithubComIfreddyrondonGocapturePayload(&r, p)
	return r.Error()
}

func (p *jsonPayload) getPayload() Payload {
	if p.Cap != nil {
		return p.Cap
	} else if p.Captures != nil {
		return p.Captures
	} else if p.Data != nil {
		return p.Data
	}
	return p.Payload
}

func (p *jsonPayload) IsValid(errors *validate.Errors) {
	payl := p.getPayload()
	if len(payl) == 0 {
		errors.Add("payload", errMissingPayload)
	}
}
