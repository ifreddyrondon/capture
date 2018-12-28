package decoder

import (
	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/capture/payload"
)

const errMissingPayload = "payload value must not be blank"

type PostPayload struct {
	Data    []payload.Metric `json:"data"`
	Payload []payload.Metric `json:"payload"`
}

func (p *PostPayload) OK() error {
	e := validate.NewErrors()
	data := getPayload(p)
	if len(data) == 0 {
		e.Add("payload", errMissingPayload)
	}

	if e.HasAny() {
		return e
	}

	p.Payload = data
	return nil
}

func (p *PostPayload) GetPayload() payload.Payload {
	return p.Payload
}

func getPayload(p *PostPayload) []payload.Metric {
	if p.Data != nil {
		return p.Data
	}
	return p.Payload
}
