package decoder

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ifreddyrondon/capture/decoder"
)

type TimestampReqDecoder interface {
	GetTimestamp() time.Time
}

// Decode gets a request payload and validate it.
func Decode(r *http.Request, v decoder.OK) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("cannot unmarshal json into timestamp value")
	}
	return v.OK()
}
