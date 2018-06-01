package payload

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/markbates/validate"
)

// errMissingPayload expected error when payload is missing
const errMissingPayload = "payload value must not be blank"

type payloadJSON struct {
	Cap      []*Metric `json:"cap"`
	Captures []*Metric `json:"captures"`
	Data     []*Metric `json:"data"`
	Payload  []*Metric `json:"payload"`
}

// UnmarshalJSON supports json.Unmarshaler interface
func (p *payloadJSON) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ad23cceDecodeGithubComIfreddyrondonCapturePayload(&r, p)
	return r.Error()
}

func (p *payloadJSON) getPayload() Payload {
	if p.Cap != nil {
		return p.Cap
	} else if p.Captures != nil {
		return p.Captures
	} else if p.Data != nil {
		return p.Data
	}
	return p.Payload
}

func (p *payloadJSON) IsValid(errors *validate.Errors) {
	payl := p.getPayload()
	if len(payl) == 0 {
		errors.Add("payload", errMissingPayload)
	}
}
