package capture

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

type payload domain.Payload

func (p payload) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *payload) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed when scan payload")
	}

	var payl payload
	if err := json.Unmarshal(source, &payl); err != nil {
		return err
	}

	*p = payl
	return nil
}
