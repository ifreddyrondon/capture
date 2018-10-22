package geocoding

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	LAT       *float64 `json:"lat"`
	LNG       *float64 `json:"lng"`
	Elevation *float64 `json:"elevation,omitempty"`
}

// Value convert Payload to a driver database Value.
func (p Point) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan assigns a value from a database driver
func (p *Point) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	return json.Unmarshal(source, p)
}
