package decoder

import (
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/capture/geocoding"
	pointDecoder "github.com/ifreddyrondon/capture/features/capture/geocoding/decoder"
	payloadDecoder "github.com/ifreddyrondon/capture/features/capture/payload/decoder"
	tagsDecoder "github.com/ifreddyrondon/capture/features/capture/tags/decoder"
	timestampDecoder "github.com/ifreddyrondon/capture/features/capture/timestamp/decoder"
	"github.com/markbates/validate"
	"gopkg.in/src-d/go-kallax.v1"
)

type POSTCapture struct {
	payloadDecoder.PostPayload
	timestampDecoder.PostTimestamp
	tagsDecoder.PostTags
	Location *pointDecoder.PostPoint `json:"location"`
}

func (c *POSTCapture) OK() error {
	e := validate.NewErrors()

	if err := c.PostPayload.OK(); err != nil {
		e.Add("payload", err.Error())
	}
	if c.Location != nil {
		if err := c.Location.OK(); err != nil {
			e.Add("location", err.Error())
		}
	}

	c.PostTimestamp.OK()

	if e.HasAny() {
		return e
	}

	return nil
}

func (c *POSTCapture) GetCapture() features.Capture {
	now := time.Now()
	var p *geocoding.Point
	if c.Location != nil {
		po := c.Location.GetPoint()
		p = &po
	}

	return features.Capture{
		ID:        kallax.NewULID(),
		Payload:   c.GetPayload(),
		Timestamp: c.GetTimestamp(),
		Tags:      c.GetTags(),
		Location:  p,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
