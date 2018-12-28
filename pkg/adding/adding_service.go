package adding

import (
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

// Store provides access to the repository storage.
type Store interface {
	CreateCapture(*domain.Capture) error
}

// Service provides adding operations.
type Service interface {
	// AddCapture add a new capture to a repository
	AddCapture(*domain.Repository, Capture) (*domain.Capture, error)
}

type service struct {
	s     Store
	clock *pkg.Clock
}

// NewService creates an adding service with the necessary dependencies
func NewService(s Store) Service {
	return &service{s: s}
}

func (s *service) AddCapture(r *domain.Repository, c Capture) (*domain.Capture, error) {
	capt := s.getDomainCapture(r, c)
	if err := s.s.CreateCapture(capt); err != nil {
		return nil, errors.Wrap(err, "could not add capture")
	}
	return capt, nil
}

func (s *service) getDomainCapture(r *domain.Repository, c Capture) *domain.Capture {
	now := time.Now()
	result := &domain.Capture{
		ID:        kallax.NewULID(),
		Payload:   c.Payload.Payload,
		Tags:      c.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if c.Location != nil {
		result.Location = &domain.Point{
			LAT:       c.Location.LAT,
			LNG:       c.Location.LNG,
			Elevation: c.Location.Elevation,
		}
	}
	if result.Tags == nil {
		result.Tags = []string{}
	}
	if c.postTimestamp != nil {
		result.Timestamp = c.postTimestamp.UTC()
	} else {
		result.Timestamp = s.clock.Now()
	}
	return result
}
