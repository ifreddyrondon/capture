package capture

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

type point domain.Point

func (p point) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *point) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed when scan point")
	}

	return json.Unmarshal(source, p)
}
