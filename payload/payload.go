package payload

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	// ErrorUnmarshalPayload expected error when fails to unmarshal a payload
	ErrorUnmarshalPayload = errors.New("cannot unmarshal json into valid payload value")
	// ErrorMissingPayload expected error when payload is missing
	ErrorMissingPayload = errors.New("missing payload value")
)

// Payload represent an association of values
type Payload map[string]interface{}

// UnmarshalJSON supports json.Unmarshaler interface
func (p *Payload) UnmarshalJSON(data []byte) error {
	var model jsonPayload
	if err := model.unmarshalJSON(data); err != nil {
		return ErrorUnmarshalPayload
	}
	payl := model.getPayload()
	if len(payl) == 0 {
		return ErrorMissingPayload
	}
	*p = payl
	return nil
}

func (v *jsonPayload) getPayload() Payload {
	if v.Cap != nil {
		return v.Cap
	} else if v.Captures != nil {
		return v.Captures
	} else if v.Data != nil {
		return v.Data
	}
	return v.Payload
}

func (p Payload) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *Payload) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed.")
	}

	var i interface{}
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed.")
	}

	return nil
}
