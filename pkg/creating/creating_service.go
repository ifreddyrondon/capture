package creating

import (
	"time"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

const defaultCrrBranchFieldValue = "master"

// Store provides access to the repository storage.
type Store interface {
	SaveRepo(*domain.Repository) error
}

// Service provides creating operations.
type Service interface {
	// CreateRepo creates a new repository to an user
	CreateRepo(*domain.User, Payload) (*Repository, error)
}

type service struct {
	s Store
}

// NewService creates a creating service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) CreateRepo(owner *domain.User, p Payload) (*Repository, error) {
	r := getDomainRepository(owner, p)
	if err := s.s.SaveRepo(r); err != nil {
		return nil, errors.Wrap(err, "could not save repo")
	}
	return getRepo(*r), nil
}

func getDomainRepository(owner *domain.User, r Payload) *domain.Repository {
	now := time.Now()
	repo := &domain.Repository{
		ID:            kallax.NewULID(),
		Name:          *r.Name,
		CurrentBranch: defaultCrrBranchFieldValue,
		CreatedAt:     now,
		UpdatedAt:     now,
		UserID:        owner.ID,
	}
	if r.Visibility == nil {
		repo.Visibility = domain.Public
	} else {
		repo.Visibility = domain.Visibility(*r.Visibility)
	}

	return repo
}

func getRepo(r domain.Repository) *Repository {
	return &Repository{
		ID:         r.ID.String(),
		Name:       r.Name,
		Visibility: string(r.Visibility),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}
