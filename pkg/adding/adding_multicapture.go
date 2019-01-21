package adding

import (
	"fmt"

	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

const (
	MultiCaptureValidator    validator.StringValidator = "cannot unmarshal json into valid multi capture value"
	errMissingCaptures                                 = "captures value must not be blank or empty"
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

	for i, capt := range m.Captures {
		if err := capt.OK(); err != nil {
			if !m.IgnoreErrors {
				key := fmt.Sprintf("capture %v", i)
				e.Add(key, fmt.Sprintf("%v: %v", key, err))
			}
		} else {
			m.CapturesOK = append(m.CapturesOK, capt)
		}
	}

	if e.HasAny() {
		return e
	}

	return nil
}
