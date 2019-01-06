package capture

import (
	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

const errMissingPayload = "payload value must not be blank"

// PayloadValidator for adding request payload
const PayloadValidator validator.StringValidator = "cannot unmarshal json into valid payload value"

type Payload struct {
	Data    []domain.Metric `json:"data"`
	Payload []domain.Metric `json:"payload"`
}

func (p *Payload) OK() error {
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
