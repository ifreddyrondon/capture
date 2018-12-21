package getting

import (
	"fmt"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type notAuthorizedErr string

func (i notAuthorizedErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i notAuthorizedErr) IsNotAuthorized() bool { return true }

// Store provides access to the repository storage.
type Store interface {
	// Get retrieve a repository from storage.
	Get(kallax.ULID) (*domain.Repository, error)
}

// Service provides listing repository operations.
type Service interface {
	// GetRepo retrieve a user repo or public repository .
	GetRepo(kallax.ULID, *domain.User) (*domain.Repository, error)
}

type service struct {
	s Store
}

// NewService creates a listing service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) GetRepo(id kallax.ULID, user *domain.User) (*domain.Repository, error) {
	repo, err := s.s.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get repo")
	}

	if repo.Visibility != domain.Public && repo.UserID != user.ID {
		errStr := fmt.Sprintf("user %v not authorized to get repo %v", user.ID, id)
		return nil, errors.WithStack(notAuthorizedErr(errStr))
	}

	return repo, nil
}
