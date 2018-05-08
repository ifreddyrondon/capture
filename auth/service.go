package auth

import (
	"errors"

	"github.com/ifreddyrondon/gocapture/user"
)

var errInvalidPassword = errors.New("invalid password")

// Service is the interface implemented by auth
// It make authentication operations.
type Service interface {
	// Authenticate validate users credentials
	Authenticate(*BasicAuthCrendential) (*user.User, error)
}

// PGAuthService implementation of auth.Service for Postgres database.
type PGAuthService struct {
	user.GetterService
}

// Authenticate validate users credentials
func (p *PGAuthService) Authenticate(crendetials *BasicAuthCrendential) (*user.User, error) {
	u, err := p.Get(crendetials.Email)
	if err != nil {
		return nil, err
	}
	if !u.CheckPassword(crendetials.Password) {
		return nil, errInvalidPassword
	}

	return u, nil
}
