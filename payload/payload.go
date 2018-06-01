package payload

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/markbates/validate"
)

var (
	// errUnmarshalPayload expected error when fails to unmarshal a payload
	errUnmarshalPayload = errors.New("cannot unmarshal json into valid payload value")
)

// Payload represent an association of metrics
type Payload []*Metric

// UnmarshalJSON supports json.Unmarshaler interface
func (p *Payload) UnmarshalJSON(data []byte) error {
	var model payloadJSON
	if err := model.unmarshalJSON(data); err != nil {
		return errUnmarshalPayload
	}

	if err := validate.Validate(&model); err.HasAny() {
		return err
	}
	*p = model.getPayload()
	return nil
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
