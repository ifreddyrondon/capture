package payload

import (
	"errors"
)

var (
	ErrorUnmarshalPayload = errors.New("cannot unmarshal json into Payload valid value")
)

// Values returns the values associated with this payload, or nil
// if no values.
type Payload interface {
	Values() interface{}
	UnmarshalJSON(data []byte) error
}
