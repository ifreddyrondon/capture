package signup

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
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
	SaveUser(user *domain.User) error
}

// Service provides sign-up operations.
type Service interface {
	// EnrollUser register a new user
	EnrollUser(Payload) (*User, error)
}

type service struct {
	s Store
}

// NewService creates an sign-up service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) EnrollUser(p Payload) (*User, error) {
	u, err := getDomainUser(p)
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
	return getUser(*u), nil
}

func getDomainUser(p Payload) (*domain.User, error) {
	now := time.Now()
	u := &domain.User{
		ID:        kallax.NewULID(),
		Email:     *p.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if p.Password == nil {
		pass, err := password.Generate(64, 10, 10, false, false)
		if err != nil {
			return nil, err
		}
		p.Password = &pass
	}

	hash, err := hashPassword(*p.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hash

	return u, nil
}

func getUser(u domain.User) *User {
	return &User{
		ID:        u.ID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func hashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), 10)
}
