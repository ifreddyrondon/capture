package creating

import (
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

const defaultCrrBranchFieldValue = "master"

// Store provides access to the repository storage.
type Store interface {
	SaveRepo(*pkg.Repository) error
}

// Service provides repository operations.
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

func getDomainRepository(owner *domain.User, r Payload) *pkg.Repository {
	// FIXME: handler id err
	ownerID, _ := kallax.NewULIDFromText(owner.ID)
	now := time.Now()
	repo := &pkg.Repository{
		ID:            kallax.NewULID(),
		Name:          *r.Name,
		CurrentBranch: defaultCrrBranchFieldValue,
		CreatedAt:     now,
		UpdatedAt:     now,
		UserID:        ownerID,
	}
	if r.Visibility == nil {
		repo.Visibility = pkg.Public
	} else {
		repo.Visibility = pkg.Visibility(*r.Visibility)
	}

	return repo
}

func getRepo(r pkg.Repository) *Repository {
	return &Repository{
		ID:         r.ID.String(),
		Name:       r.Name,
		Visibility: string(r.Visibility),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}
