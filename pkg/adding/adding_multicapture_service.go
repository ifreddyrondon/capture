package adding

import (
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// CapturesStore provides access to the captures storage.
type CapturesStore interface {
	CreateCaptures([]domain.Capture) error
}

// CaptureService provides adding operations.
type CapturesService interface {
	// AddCaptures add new captures to a repository
	AddCaptures(*domain.Repository, MultiCapture) ([]domain.Capture, error)
}

type capturesService struct {
	s CapturesStore
}

// NewCapturesService creates an adding service with the necessary dependencies to add captures.
func NewCapturesService(s CapturesStore) CapturesService {
	return &capturesService{s: s}
}

func (s *capturesService) AddCaptures(r *domain.Repository, captures MultiCapture) ([]domain.Capture, error) {
	capt := getDomainCaptures(r, captures)
	if err := s.s.CreateCaptures(capt); err != nil {
		return nil, errors.Wrap(err, "could not add captures")
	}
	return capt, nil
}

func getDomainCaptures(r *domain.Repository, m MultiCapture) []domain.Capture {
	// TODO: make real captures
	result := make([]domain.Capture, len(m.Captures))
	return result
}
