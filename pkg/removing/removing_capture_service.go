package removing

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/domain"
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
		errStr := fmt.Sprintf("could not remove capture %v", c.ID)
		return errors.Wrap(err, errStr)
	}
	return nil
}
