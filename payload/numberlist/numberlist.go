package numberlist

import (
	"errors"

	"log"
)

var (
	// ErrorUnmarshalPayload expected error when fails to unmarshal a payload
	ErrorUnmarshalPayload = errors.New("cannot unmarshal json into Payload valid value")
)

// Payload represent an association of float numbers
type Payload []float64

// New returns a new pointer to a ArrayNumberPayload composed of the passed float64
func New(data ...float64) *Payload {
	p := new(Payload)
	*p = data
	return p
}

// UnmarshalJSON supports json.Unmarshaler interface
func (p *Payload) UnmarshalJSON(data []byte) error {
	var model jsonPayload
	if err := model.unmarshalJSON(data); err != nil {
		log.Print(err)
		return ErrorUnmarshalPayload
	}
	*p = model.getPayload()
	return nil
}

func (v *jsonPayload) getPayload() []float64 {
	if v.Cap != nil {
		return v.Cap
	} else if v.Captures != nil {
		return v.Captures
	} else if v.Data != nil {
		return v.Data
	}
	return v.Payload
}
