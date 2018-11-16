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

func (s *Service) GetUserRepositories(u *features.User, l *listing.Listing) ([]features.Repository, error) {
	listingRepo := newListingRepo(*l)
	listingRepo.Owner = u
	return s.Store.List(listingRepo)
}

func (s *Service) GetPublicRepositories(l *listing.Listing) ([]features.Repository, error) {
	listingRepo := newListingRepo(*l)
	listingRepo.Visibility = &features.Public
	return s.Store.List(listingRepo)
}

func (s *Service) GetRepo(id kallax.ULID, loggedUser *features.User) (*features.Repository, error) {
	repo, err := s.Store.Get(id)
	if err != nil {
		return nil, err
	}
	if repo.Visibility != features.Public && repo.UserID != loggedUser.ID {
		return nil, ErrorNotAuthorized
	}
	return repo, nil
}
