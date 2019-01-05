package removing

import (
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// CaptureStore provides access to the capture storage.
type CaptureStore interface {
	// Delete a capture from a repo.
	Delete(*domain.Capture) error
}

// CaptureService provides removing capture operations.
type CaptureService interface {
	// Remove a repo capture from a repo.
	Remove(*domain.Capture) error
}

type captureService struct {
	s CaptureStore
}

// NewCaptureService creates a getting service with the necessary dependencies
func NewCaptureService(s CaptureStore) CaptureService {
	return &captureService{s: s}
}

func (s *captureService) Remove(c *domain.Capture) error {
	err := s.s.Delete(c)
	if err != nil {
		return errors.Wrap(err, "could not remove capture")
	}
	return nil
}
