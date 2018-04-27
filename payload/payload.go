package payload

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrorUnmarshalPayload expected error when fails to unmarshal a payload
	ErrorUnmarshalPayload = errors.New("cannot unmarshal json into valid payload value")
	// ErrorMissingPayload expected error when payload is missing
	ErrorMissingPayload = errors.New("missing payload value")
)

// Payload represent an association of metrics
type Payload []*Metric

// UnmarshalJSON supports json.Unmarshaler interface
func (p *Payload) UnmarshalJSON(data []byte) error {
	var model jsonPayload
	if err := model.unmarshalJSON(data); err != nil {
		fmt.Println(err)
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

// Value convert Payload to a driver database Value.
func (p Payload) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan assigns a value from a database driver
func (p *Payload) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var payl []*Metric
	if err := json.Unmarshal(source, &payl); err != nil {
		return err
	}

	*p = payl

	return nil
}
