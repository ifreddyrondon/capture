package adding

import (
	"github.com/gobuffalo/validate"

	"github.com/ifreddyrondon/capture/pkg/validator"

	"github.com/pkg/errors"
)

type Capture struct {
	validator.Payload
	validator.Timestamp
	Tags     []string               `json:"tags"`
	Location *validator.GeoLocation `json:"location"`
}

func (c *Capture) Validate() error {
	e := validate.NewErrors()
	if err := c.Payload.Validate(); err != nil {
		e.Add("payload", err.Error())
	}
	if err := c.Timestamp.Validate(); err != nil {
		err = errors.Wrap(err, "invalid timestamp value")
		e.Add("timestamp", err.Error())
	}
	if c.Location != nil {
		if err := c.Location.Validate(); err != nil {
			e.Add("location", err.Error())
		}
	}
	if e.HasAny() {
		return e
	}

	return nil
}
