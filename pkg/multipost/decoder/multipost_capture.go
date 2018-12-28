package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/gobuffalo/validate"
)

const (
	maxAllowedItemsToPost = 10
)

var errMaxAllowedItemsToPost string

func init() {
	errMaxAllowedItemsToPost = fmt.Sprintf("the maximum amount of allowed captures is %v", maxAllowedItemsToPost)
}

type MultiPOSTCaptures struct {
	IgnoreErrors  bool `json:"ignore_errors"`
	Notifications struct {
		CallbackURL string `json:"callback_url"`
		Email       string `json:"email"`
	}
	Captures []json.RawMessage `json:"captures"`
}

func (c *MultiPOSTCaptures) OK() error {
	e := validate.NewErrors()
	if !c.IgnoreErrors && len(c.Captures) > maxAllowedItemsToPost {
		e.Add("captures", errMaxAllowedItemsToPost)
	}

	if e.HasAny() {
		return e
	}

	return nil
}
