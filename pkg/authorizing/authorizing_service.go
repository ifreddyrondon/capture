package authorizing

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

type invalidIDErr string

func (i invalidIDErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidIDErr) IsInvalid() bool { return true }

type invalidCredentialErr string

func (i invalidCredentialErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i invalidCredentialErr) IsNotAuthorized() bool { return true }

type notFoundErr interface {
	NotFound() bool
}

func isNotFound(err error) bool {
	if e, ok := errors.Cause(err).(notFoundErr); ok {
		return e.NotFound()
	}
	return false
}

// Store provides access to the user storage.
type Store interface {
	// GetUserByID get user by id.
	GetUserByID(kallax.ULID) (*domain.User, error)
}

// TokenService provides utils to handle authorizing token.
type TokenService interface {
	// IsRequestAuthorized validates if a request if authorized
	IsRequestAuthorized(*http.Request) (string, error)
}

// Service provides authorizing operations.
type Service interface {
	AuthorizeRequest(*http.Request) (*domain.User, error)
}

type service struct {
	s  Store
	ts TokenService
}

// NewService creates an authenticating service with the necessary dependencies
func NewService(ts TokenService, s Store) Service {
	return &service{ts: ts, s: s}
}

func (s *service) AuthorizeRequest(r *http.Request) (*domain.User, error) {
	subjectID, err := s.ts.IsRequestAuthorized(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not authorized request")
	}
	id, err := kallax.NewULIDFromText(subjectID)
	if err != nil {
		return nil, invalidIDErr(fmt.Sprintf("%v is not a valid ULID", subjectID))
	}
	u, err := s.s.GetUserByID(id)
	if err != nil {
		if isNotFound(err) {
			return nil, errors.WithStack(invalidCredentialErr(err.Error()))
		}
		return nil, errors.Wrap(err, "error when get user by id in AuthorizeRequest")
	}
	return u, nil
}
