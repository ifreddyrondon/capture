package updating

import (
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/validator"
)

type Capture struct {
	*validator.Payload
	*validator.Timestamp
	Location *validator.GeoLocation `json:"location"`
	Tags     []string               `json:"tags"`
}

func (c *Capture) Validate() error {
	e := validate.NewErrors()
	if c.Payload != nil {
		if err := c.Payload.Validate(); err != nil {
			e.Add("payload", err.Error())
		}
	}
	if c.Timestamp != nil {
		if err := c.Timestamp.Validate(); err != nil {
			err = errors.Wrap(err, "invalid timestamp value")
			e.Add("timestamp", err.Error())
		}
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
