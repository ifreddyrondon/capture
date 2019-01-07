package getting

import (
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

// RepoStore provides access to the repository storage.
type RepoStore interface {
	// Get retrieve a repository from storage.
	Get(kallax.ULID) (*domain.Repository, error)
}

// RepoService provides getting repository operations.
type RepoService interface {
	// Get retrieve a repo.
	Get(kallax.ULID) (*domain.Repository, error)
}

type repoService struct {
	s RepoStore
}

// NewRepoService creates a listing service with the necessary dependencies
func NewRepoService(s RepoStore) RepoService {
	return &repoService{s: s}
}

func (s *repoService) Get(id kallax.ULID) (*domain.Repository, error) {
	repo, err := s.s.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get repo")
	}
	return repo, nil
}
