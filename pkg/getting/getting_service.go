package getting

import (
	"fmt"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type notAuthorizedErr string

func (i notAuthorizedErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i notAuthorizedErr) IsNotAuthorized() bool { return true }

// Store provides access to the repository storage.
type Store interface {
	// Get retrieve a repository from storage.
	Get(string) (*pkg.Repository, error)
}

// Service provides listing repository operations.
type Service interface {
	// GetRepo retrieve a user repo or public repository .
	GetRepo(string, *pkg.User) (*pkg.Repository, error)
}

type service struct {
	s Store
}

// NewService creates a listing service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) GetRepo(id string, user *pkg.User) (*pkg.Repository, error) {
	repo, err := s.s.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get repo")
	}

	loggedUsrID, _ := kallax.NewULIDFromText(user.ID)
	if repo.Visibility != pkg.Public && repo.UserID != loggedUsrID {
		errStr := fmt.Sprintf("user %v not authorized to get repo %v", user.ID, id)
		return nil, errors.WithStack(notAuthorizedErr(errStr))
	}

	return repo, nil
}
