package decoder

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/capture/geocoding"
	"github.com/markbates/validate"
	"gopkg.in/src-d/go-kallax.v1"
)

const errIDMissing = "capture id must not be blank"

type PUTCapture struct {
	ID *kallax.ULID `json:"id"`
	POSTCapture
}

func (c *PUTCapture) OK() error {
	if err := c.POSTCapture.OK(); err != nil {
		return err
	}

	e := validate.NewErrors()
	if c.ID == nil {
		e.Add("id", errIDMissing)
	}

	if e.HasAny() {
		return e
	}

	return nil
}

func (c *PUTCapture) GetCapture() pkg.Capture {
	var p *geocoding.Point
	if c.Location != nil {
		po := c.Location.GetPoint()
		p = &po
	}

	return pkg.Capture{
		ID:        *c.ID,
		Payload:   c.GetPayload(),
		Timestamp: c.GetTimestamp(),
		Tags:      c.GetTags(),
		Location:  p,
	}
}
