package adding

import (
	"fmt"

	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

const (
	MultiCaptureValidator    validator.StringValidator = "cannot unmarshal json into valid multi capture value"
	errMissingCaptures                                 = "captures value must not be blank"
	maxAllowedCapturesToPost                           = 50
)

type MultiCapture struct {
	IgnoreErrors bool      `json:"ignore_errors"`
	Captures     []Capture `json:"captures"`
	CapturesOK   []Capture `json:"-"`
}

func (m *MultiCapture) OK() error {
	e := validate.NewErrors()
	if len(m.Captures) > maxAllowedCapturesToPost {
		e.Add("captures", fmt.Sprintf("the maximum amount of allowed captures is %v", maxAllowedCapturesToPost))
		return e
	}

	if len(m.Captures) == 0 {
		e.Add("captures", errMissingCaptures)
		return e
	}

	for i, c := range m.Captures {
		err := c.OK()
		if err != nil && !m.IgnoreErrors {
			e.Add(fmt.Sprintf("capture %v", i), err.Error())
		} else {
			m.CapturesOK = append(m.CapturesOK, c)
		}
	}

	if e.HasAny() {
		return e
	}

	return nil
}
