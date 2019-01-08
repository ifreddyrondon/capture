package listing

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// RepoStore provides access to the repository storage.
type RepoStore interface {
	// List retrieve repositories with domain.Listing attrs.
	List(*domain.Listing) ([]domain.Repository, int64, error)
}

// RepoService provides listing repository operations.
type RepoService interface {
	// GetUserRepos get the repositories from a given user.
	GetUserRepos(*domain.User, *listing.Listing) (*ListRepositoryResponse, error)
	// GetPublicRepos get all the public repos.
	GetPublicRepos(*listing.Listing) (*ListRepositoryResponse, error)
}

type repoService struct {
	s RepoStore
}

// NewRepoService creates a listing service with the necessary dependencies
func NewRepoService(s RepoStore) RepoService {
	return &repoService{s: s}
}

func (s *repoService) GetUserRepos(u *domain.User, l *listing.Listing) (*ListRepositoryResponse, error) {
	lrepo := domain.NewListing(*l)
	lrepo.Owner = &u.ID
	repos, total, err := s.s.List(lrepo)
	if err != nil {
		return nil, errors.Wrap(err, "err getting user repos")
	}
	l.Paging.Total = total
	return newListRepoResponse(repos, l), nil
}

func (s *repoService) GetPublicRepos(l *listing.Listing) (*ListRepositoryResponse, error) {
	lrepo := domain.NewListing(*l)
	lrepo.Visibility = &domain.Public
	repos, total, err := s.s.List(lrepo)
	if err != nil {
		return nil, errors.Wrap(err, "err getting public repos")
	}
	l.Paging.Total = total
	return newListRepoResponse(repos, l), nil
}

type ListRepositoryResponse struct {
	Results []domain.Repository `json:"results"`
	Listing *listing.Listing    `json:"listing"`
}

func newListRepoResponse(repos []domain.Repository, l *listing.Listing) *ListRepositoryResponse {
	if repos == nil {
		repos = make([]domain.Repository, 0)
	}
	return &ListRepositoryResponse{Listing: l, Results: repos}
}
