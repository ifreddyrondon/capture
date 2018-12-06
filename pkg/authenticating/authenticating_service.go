package authenticating

import (
	"fmt"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type notFoundErr interface {
	// NotFound returns true when a resource is not found.
	NotFound() bool
}

func isNotFound(err error) bool {
	if e, ok := errors.Cause(err).(notFoundErr); ok {
		return e.NotFound()
	}
	return false
}

type invalidCredentialErr string

func (i invalidCredentialErr) Error() string            { return fmt.Sprintf(string(i)) }
func (i invalidCredentialErr) InvalidCredentials() bool { return true }

const errInvalidPassword invalidCredentialErr = "invalid password"

// Store provides access to the user storage.
type Store interface {
	// GetUserByEmail get user by email.
	GetUserByEmail(string) (*pkg.User, error)
}

// TokenService provides utils to handle authentication token.
type TokenService interface {
	// GenerateToken an authorization token.
	GenerateToken(string) (string, error)
}

// Service provides authenticating operations.
type Service interface {
	// AuthenticateUser compare the given credentials with the stored ones.
	AuthenticateUser(BasicCredential) (*pkg.User, error)
	GetUserToken(string) (string, error)
}

type service struct {
	s  Store
	ts TokenService
}

// NewService creates an authenticating service with the necessary dependencies
func NewService(ts TokenService, s Store) Service {
	return &service{ts: ts, s: s}
}

// GenerateToken creates a new token
func (s *service) GetUserToken(userID string) (string, error) {
	t, err := s.ts.GenerateToken(userID)
	if err != nil {
		return "", errors.Wrap(err, "could not generate token")
	}
	return t, nil
}

func (s *service) AuthenticateUser(credential BasicCredential) (*pkg.User, error) {
	u, err := s.s.GetUserByEmail(credential.Email)
	if err != nil {
		if isNotFound(err) {
			return nil, errors.WithStack(invalidCredentialErr(err.Error()))
		}
		return nil, err
	}

	if !checkPassword(u.Password, []byte(credential.Password)) {
		return nil, errors.WithStack(errInvalidPassword)
	}

	return u, nil
}

func checkPassword(hashed, pass []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashed, pass); err != nil {
		return false
	}
	return true
}
