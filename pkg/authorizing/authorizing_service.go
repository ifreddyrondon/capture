package authorizing

import (
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type invalidCredentialErr string

func (i invalidCredentialErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i invalidCredentialErr) IsNotAuthorized() bool { return true }

// Store provides access to the user storage.
type Store interface {
	// GetUserByEmail get user by email.
	GetUserByID(kallax.ULID) (*pkg.User, error)
}

// TokenService provides utils to handle authorizing token.
type TokenService interface {
	// IsRequestAuthorized validates if a request if authorized
	IsRequestAuthorized(*http.Request) (string, error)
}

// Service provides authorizing operations.
type Service interface {
	AuthorizeRequest(*http.Request) (*pkg.User, error)
}

type service struct {
	s  Store
	ts TokenService
}

// NewService creates an authenticating service with the necessary dependencies
func NewService(ts TokenService, s Store) Service {
	return &service{ts: ts, s: s}
}

func (s *service) AuthorizeRequest(r *http.Request) (*pkg.User, error) {
	subjectID, err := s.ts.IsRequestAuthorized(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not authorized request")
	}

	userID, err := kallax.NewULIDFromText(subjectID)
	if err != nil {
		return nil, errors.WithStack(invalidCredentialErr(err.Error()))
	}
	return s.s.GetUserByID(userID)
}
