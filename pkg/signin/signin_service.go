package signin

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/pkg/errors"
)

// Store provides access to the user storage.
type Store interface {
	SaveUser(user *pkg.User) error
}

// Service provides authenticating operations.
type Service interface {
	// EnrollUser register a new user
	EnrollUser(*pkg.User) error
}

type service struct {
	s Store
}

// NewService creates an signin service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) EnrollUser(u *pkg.User) error {
	err := s.s.SaveUser(u)
	if err != nil {
		return errors.Wrap(err, "EnrollUser")
	}
	return nil
}
