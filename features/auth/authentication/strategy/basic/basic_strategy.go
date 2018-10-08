package basic

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/user"
)

var errInvalidCredentials = errors.New("invalid email or password")

type decodingErr struct {
	err error
}

func (d *decodingErr) Error() string {
	return d.err.Error()
}

// Basic strategy mechanisms
type Basic struct {
	service user.GetterService
}

// New returns a new instance of basic strategy
func New(service user.GetterService) *Basic {
	return &Basic{service: service}
}

// Validate basic credentials.
func (b *Basic) Validate(r *http.Request) (*features.User, error) {
	var cre Credentials
	if err := json.NewDecoder(r.Body).Decode(&cre); err != nil {
		return nil, &decodingErr{err: err}
	}

	u, err := b.service.GetByEmail(cre.Email)
	if err != nil {
		if err == user.ErrNotFound {
			return nil, errInvalidCredentials
		}
		return nil, err
	}
	if !u.CheckPassword(cre.Password) {
		return nil, errInvalidCredentials
	}

	return u, nil
}

// IsErrCredentials check if an error is for invalid credentials.
func (b *Basic) IsErrCredentials(err error) bool {
	return err == errInvalidCredentials
}

// IsErrDecoding check if an error is for invalid decoding credentials.
func (b *Basic) IsErrDecoding(err error) bool {
	if _, ok := err.(*decodingErr); ok {
		return true
	}
	return false
}
