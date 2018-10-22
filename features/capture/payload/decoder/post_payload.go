package decoder

import (
	"github.com/ifreddyrondon/capture/features/capture/payload"
	"github.com/markbates/validate"
)

const errMissingPayload = "payload value must not be blank"

type PostPayload struct {
	Cap      []payload.Metric `json:"cap"`
	Captures []payload.Metric `json:"captures"`
	Data     []payload.Metric `json:"data"`
	Payload  []payload.Metric `json:"payload"`
}

func (p *PostPayload) OK() error {
	e := validate.NewErrors()
	captures := getPayload(p)
	if len(captures) == 0 {
		e.Add("payload", errMissingPayload)
	}

	if e.HasAny() {
		return e
	}

	p.Payload = captures
	return nil
}

func (p *PostPayload) GetPayload() payload.Payload {
	return p.Payload
}

func getPayload(p *PostPayload) []payload.Metric {
	if p.Cap != nil {
		return p.Cap
	} else if p.Captures != nil {
		return p.Captures
	} else if p.Data != nil {
		return p.Data
	}
	return p.Payload
}
