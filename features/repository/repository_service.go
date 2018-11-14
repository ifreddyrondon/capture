package repository

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/features"
	"gopkg.in/src-d/go-kallax.v1"
)

type Service struct {
	Store
}

func (s *Service) Save(u *features.User, r *features.Repository) error {
	return s.Store.Save(u, r)
}

func (s *Service) GetUserRepos(u *features.User, l *listing.Listing) ([]features.Repository, error) {
	listingRepo := newListingRepo(*l)
	listingRepo.Owner = u
	return s.Store.List(listingRepo)
}

func (s *Service) GetPublicRepos(l *listing.Listing) ([]features.Repository, error) {
	listingRepo := newListingRepo(*l)
	listingRepo.Visibility = &features.Public
	return s.Store.List(listingRepo)
}

func (s *Service) GetUserRepo(u *features.User, id kallax.ULID) (*features.Repository, error) {
	repo, err := s.Store.Get(id)
	if err != nil {
		return nil, err
	}
	if repo.UserID != u.ID {
		return nil, ErrorNotAuthorized
	}
	return repo, nil
}
