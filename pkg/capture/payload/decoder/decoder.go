package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/decoder"
	"github.com/ifreddyrondon/capture/pkg/capture/payload"
)

type PayloadReqDecoder interface {
	GetPayload() payload.Payload
}

// Decode gets a request payload and validate it.
func Decode(r *http.Request, v decoder.OK) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("cannot unmarshal json into valid payload value")
	}
	return v.OK()
}
