package signin

import (
	"fmt"
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/src-d/go-kallax.v1"
)

type conflictErr string

func (e conflictErr) Error() string  { return string(e) }
func (e conflictErr) Conflict() bool { return true }

type constraintErr interface {
	UniqueConstraint() bool
}

func isConstraintErr(err error) bool {
	if err, ok := errors.Cause(err).(constraintErr); ok {
		return err.UniqueConstraint()
	}
	return false
}

// Store provides access to the user storage.
type Store interface {
	SaveUser(user *pkg.User) error
}

// Service provides authenticating operations.
type Service interface {
	// EnrollUser register a new user
	EnrollUser(Payload) (*pkg.User, error)
}

type service struct {
	s Store
}

// NewService creates an signin service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) EnrollUser(p Payload) (*pkg.User, error) {
	u, err := getUser(p)
	if err != nil {
		return nil, errors.Wrap(err, "could not get user from payload")
	}
	if err := s.s.SaveUser(u); err != nil {
		if isConstraintErr(err) {
			e := conflictErr(fmt.Sprintf("email %v already exist", u.Email))
			return nil, errors.WithStack(e)
		}
		return nil, errors.Wrap(err, "could not save user")
	}
	return u, nil
}

func getUser(p Payload) (*pkg.User, error) {
	now := time.Now()
	u := &pkg.User{
		ID:        kallax.NewULID(),
		Email:     *p.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if p.Password != nil {
		hash, err := hashPassword(*p.Password)
		if err != nil {
			return nil, err
		}
		u.Password = hash
	}

	return u, nil
}

func hashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), 10)
}
