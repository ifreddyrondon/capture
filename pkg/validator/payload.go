package validator

import (
	"github.com/gobuffalo/validate"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

const errMissingPayload = "payload value must not be blank"

type Payload struct {
	Data    []domain.Metric `json:"data"`
	Payload []domain.Metric `json:"payload"`
}

func (p *Payload) Validate() error {
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

func getPayload(p *Payload) []domain.Metric {
	if p.Data != nil {
		return p.Data
	}
	return p.Payload
}
