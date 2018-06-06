package basic

import (
	"errors"

	"github.com/markbates/validate"

	"github.com/ifreddyrondon/capture/app/user"
)

var errInvalidCredentials = errors.New("invalid email or password")

// Basic strategy mechanisms
type Basic struct {
	service user.GetterService
}

// New returns a new instance of basic strategy
func New(service user.GetterService) *Basic {
	return &Basic{service: service}
}

// Validate basic credentials.
func (b *Basic) Validate(payload []byte) (*user.User, error) {
	var cre Crendentials
	if err := cre.UnmarshalJSON(payload); err != nil {
		return nil, err
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
	if _, ok := err.(*validate.Errors); ok {
		return true
	}
	if err == errInvalidPayload {
		return true
	}
	return false
}
