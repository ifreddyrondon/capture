package repository

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg"
	"gopkg.in/src-d/go-kallax.v1"
)

type Service struct {
	Store
}

func (s *Service) Save(u *pkg.User, r *pkg.Repository) error {
	return s.Store.Save(u, r)
}

func (s *Service) GetUserRepositories(u *pkg.User, l *listing.Listing) ([]pkg.Repository, error) {
	listingRepo := newListingRepo(*l)
	listingRepo.Owner = u
	return s.Store.List(listingRepo)
}

func (s *Service) GetPublicRepositories(l *listing.Listing) ([]pkg.Repository, error) {
	listingRepo := newListingRepo(*l)
	listingRepo.Visibility = &pkg.Public
	return s.Store.List(listingRepo)
}

func (s *Service) GetRepo(id string, loggedUser *pkg.User) (*pkg.Repository, error) {
	repo, err := s.Store.Get(id)
	if err != nil {
		return nil, err
	}
	// FIXME: handler id err
	loggedUsrID, _ := kallax.NewULIDFromText(loggedUser.ID)
	if repo.Visibility != pkg.Public && repo.UserID != loggedUsrID {
		return nil, ErrorNotAuthorized
	}

	return repo, nil
}
