package adding

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
)

// CaptureStore provides access to the capture storage.
type CaptureStore interface {
	CreateCapture(*domain.Capture) error
}

// CaptureService provides adding operations.
type CaptureService interface {
	// AddCapture add a new capture to a repository
	AddCapture(*domain.Repository, Capture) (*domain.Capture, error)
}

type captureService struct {
	s     CaptureStore
	clock *pkg.Clock
}

// NewCaptureService creates an adding service with the necessary dependencies
func NewCaptureService(s CaptureStore) CaptureService {
	return &captureService{s: s}
}

func (s *captureService) AddCapture(r *domain.Repository, c Capture) (*domain.Capture, error) {
	capt := getDomainCapture(s.clock, r, c)
	if err := s.s.CreateCapture(capt); err != nil {
		return nil, errors.Wrap(err, "could not add capture")
	}
	return capt, nil
}

func getDomainCapture(clock *pkg.Clock, r *domain.Repository, c Capture) *domain.Capture {
	now := time.Now()
	result := &domain.Capture{
		ID:           kallax.NewULID(),
		Payload:      c.Payload.Payload,
		Tags:         c.Tags,
		CreatedAt:    now,
		UpdatedAt:    now,
		RepositoryID: r.ID,
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
	if c.Timestamp.Time != nil {
		result.Timestamp = c.Timestamp.Time.UTC()
	} else {
		result.Timestamp = clock.Now()
	}
	return result
}
