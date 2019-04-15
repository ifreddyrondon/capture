package updating

import (
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/validator"
)

const CaptureValidator validator.StringValidator = "cannot unmarshal json into valid capture value"

type Capture struct {
	*validator.Payload
	*validator.Timestamp
	Location *validator.GeoLocation `json:"location"`
	Tags     []string               `json:"tags"`
}

func (c *Capture) OK() error {
	e := validate.NewErrors()
	if c.Payload != nil {
		if err := c.Payload.OK(); err != nil {
			e.Add("payload", err.Error())
		}
	}
	if c.Timestamp != nil {
		if err := c.Timestamp.OK(); err != nil {
			err = errors.Wrap(err, "invalid timestamp value")
			e.Add("timestamp", err.Error())
		}
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
