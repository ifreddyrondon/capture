package removing

import (
	"time"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// CaptureStore provides access to the capture storage.
type CaptureStore interface {
	// Save the capture state into the storage.
	Save(*domain.Capture) error
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
	t := time.Now()
	c.DeletedAt = &t
	err := s.s.Save(c)
	if err != nil {
		return errors.Wrap(err, "could not remove capture")
	}
	return nil
}
