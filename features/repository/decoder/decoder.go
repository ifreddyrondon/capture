package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/features"
)

type decoder interface {
	ok() error
	repository(repository *features.Repository)
}

// Decode gets a body payload transforming
func Decode(r *http.Request, v decoder) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("cannot unmarshal json into valid repository")
	}
	return v.ok()
}

// Repository returns the repository representation of a decoder struct.
func Repository(dec decoder, repo *features.Repository) {
	dec.repository(repo)
}
