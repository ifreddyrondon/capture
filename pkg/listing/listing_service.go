package listing

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// Store provides access to the repository storage.
type Store interface {
	// List retrieve repositories with domain.Listing attrs.
	List(*domain.Listing) ([]pkg.Repository, error)
}

// Service provides listing repository operations.
type Service interface {
	// GetUserRepos get the repositories from a given user.
	GetUserRepos(*pkg.User, *listing.Listing) (*ListRepositoryResponse, error)
	// GetPublicRepos get all the public repos.
	GetPublicRepos(*listing.Listing) (*ListRepositoryResponse, error)
}

type service struct {
	s Store
}

// NewService creates a listing service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) GetUserRepos(u *pkg.User, l *listing.Listing) (*ListRepositoryResponse, error) {
	lrepo := domain.NewListing(*l)
	lrepo.Owner = u.ID
	repos, err := s.s.List(lrepo)
	if err != nil {
		return nil, errors.Wrap(err, "err getting user repos")
	}
	return &ListRepositoryResponse{Listing: l, Results: repos}, nil
}

func (s *service) GetPublicRepos(l *listing.Listing) (*ListRepositoryResponse, error) {
	lrepo := domain.NewListing(*l)
	lrepo.Visibility = &pkg.Public
	repos, err := s.s.List(lrepo)
	if err != nil {
		return nil, errors.Wrap(err, "err getting public repos")
	}
	return &ListRepositoryResponse{Listing: l, Results: repos}, nil
}

type ListRepositoryResponse struct {
	Results []pkg.Repository `json:"results"`
	Listing *listing.Listing `json:"listing"`
}
