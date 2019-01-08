package adding

import (
	"fmt"

	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

const MultiCaptureValidator validator.StringValidator = "cannot unmarshal json into valid multi capture value"

const (
	maxAllowedCapturesToPost = 10
)

type maxAllowedCapturesErr string

func (e maxAllowedCapturesErr) Error() string    { return string(e) }
func (e maxAllowedCapturesErr) MaxAllowed() bool { return true }

type MultiCapture []Capture

func (c *MultiCapture) OK() error {
	e := validate.NewErrors()
	if len(*c) > maxAllowedCapturesToPost {
		errStr := fmt.Sprintf("the maximum amount of allowed captures is %v", maxAllowedCapturesToPost)
		e.Add("captures", maxAllowedCapturesErr(errStr).Error())
	}

	if e.HasAny() {
		return e
	}

	return nil
}
