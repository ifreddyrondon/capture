package adding

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// MultiCaptureStore provides access to the captures storage.
type MultiCaptureStore interface {
	CreateCaptures(...domain.Capture) error
}

// MultiCaptureService provides adding operations.
type MultiCaptureService interface {
	// AddCaptures add new captures to a repository
	AddCaptures(*domain.Repository, MultiCapture) ([]domain.Capture, error)
}

type multiCaptureService struct {
	s     MultiCaptureStore
	clock *pkg.Clock
}

// NewMultiCaptureService creates an adding service with the necessary dependencies to add captures.
func NewMultiCaptureService(s MultiCaptureStore) MultiCaptureService {
	return &multiCaptureService{s: s}
}

func (s *multiCaptureService) AddCaptures(r *domain.Repository, multiCapture MultiCapture) ([]domain.Capture, error) {
	captures := getDomainCaptures(s.clock, r, multiCapture.CapturesOK)
	if err := s.s.CreateCaptures(captures...); err != nil {
		return nil, errors.Wrap(err, "could not add captures")
	}
	return captures, nil
}

func getDomainCaptures(clock *pkg.Clock, r *domain.Repository, captures []Capture) []domain.Capture {
	result := make([]domain.Capture, len(captures))
	for i, c := range captures {
		result[i] = *getDomainCapture(clock, r, c)
	}

	return result
}
