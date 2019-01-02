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

// RepoStore provides access to the repository storage.
type RepoStore interface {
	// Get retrieve a repository from storage.
	Get(kallax.ULID) (*domain.Repository, error)
}

// RepoService provides listing repository operations.
type RepoService interface {
	// Get retrieve a user repo or public repository .
	Get(kallax.ULID, *domain.User) (*domain.Repository, error)
}

type service struct {
	s RepoStore
}

// NewRepoService creates a listing service with the necessary dependencies
func NewRepoService(s RepoStore) RepoService {
	return &service{s: s}
}

func (s *service) Get(id kallax.ULID, user *domain.User) (*domain.Repository, error) {
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
