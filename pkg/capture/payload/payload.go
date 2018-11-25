package payload

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Payload represent an association of metrics
type Payload []Metric

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

	var payl []Metric
	if err := json.Unmarshal(source, &payl); err != nil {
		return err
	}

	*p = payl

	return nil
}
