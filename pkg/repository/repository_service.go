package repository

import (
	"github.com/ifreddyrondon/capture/pkg"
	"gopkg.in/src-d/go-kallax.v1"
)

type Service struct {
	Store
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
