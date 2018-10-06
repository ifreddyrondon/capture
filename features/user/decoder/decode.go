package decoder

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ok interface {
	ok() error
}

// Decode gets a body payload transforming
func Decode(r *http.Request, v ok) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("cannot unmarshal json into valid user")
	}
	return v.ok()
}
