package validator

import (
	"encoding/json"
	"errors"
	"net/http"
)

// OK represents types capable of validating themselves
type OK interface {
	OK() error
}

type Validator interface {
	// Valid gets inputs from request and validate it
	Valid(r *http.Request, v OK) error
}

type StringValidator string

// DefaultJSONValidator returns a default msg when unmarshal json request body fails.
const DefaultJSONValidator StringValidator = "cannot unmarshal json body"

// Decode gets a request payload and validate it.
func (s StringValidator) Decode(r *http.Request, v OK) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New(string(s))
	}
	return v.OK()
}
