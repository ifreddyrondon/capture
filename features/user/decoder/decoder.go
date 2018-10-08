package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/features"
)

type decoder interface {
	ok() error
	user(*features.User) error
}

// Decode gets a body payload transforming
func Decode(r *http.Request, v decoder) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("cannot unmarshal json into valid user")
	}
	return v.ok()
}

// User returns the user representation of a decoder struct.
func User(dec decoder, usr *features.User) error {
	return dec.user(usr)
}
