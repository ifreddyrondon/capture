package capture

import (
	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/validator"

	"github.com/pkg/errors"
)

const Validator validator.StringValidator = "cannot unmarshal json into valid capture value"

type Capture struct {
	Payload
	Timestamp
	Tags     []string     `json:"tags"`
	Location *GeoLocation `json:"location"`
}

func (c *Capture) OK() error {
	e := validate.NewErrors()
	if err := c.Payload.OK(); err != nil {
		e.Add("payload", err.Error())
	}
	if err := c.Timestamp.OK(); err != nil {
		err = errors.Wrap(err, "invalid timestamp value")
		e.Add("timestamp", err.Error())
	}
	if c.Location != nil {
		if err := c.Location.OK(); err != nil {
			e.Add("location", err.Error())
		}
	}
	if e.HasAny() {
		return e
	}

	return nil
}
