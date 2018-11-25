package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/decoder"
	"github.com/ifreddyrondon/capture/pkg/capture/geocoding"
)

type PointReqDecoder interface {
	GetPoint() geocoding.Point
}

// Decode gets a request payload and validate it.
func Decode(r *http.Request, v decoder.OK) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("cannot unmarshal json into point value")
	}
	return v.OK()
}
