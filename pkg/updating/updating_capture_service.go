package updating

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

// CaptureService provides updating capture operations.
type CaptureService interface {
	// Update a repo capture.
	Update(Capture, *domain.Capture) error
}

type captureService struct {
	s CaptureStore
}

// NewCaptureService creates a getting service with the necessary dependencies
func NewCaptureService(s CaptureStore) CaptureService {
	return &captureService{s: s}
}

func (s *captureService) Update(data Capture, c *domain.Capture) error {
	updateCapture(data, c)
	if err := s.s.Save(c); err != nil {
		errStr := fmt.Sprintf("could not update capture %v", c.ID)
		return errors.Wrap(err, errStr)
	}
	return nil
}

func updateCapture(data Capture, capt *domain.Capture) {
	capt.UpdatedAt = time.Now()
	if data.Payload != nil {
		capt.Payload = data.Payload.Payload
	}
	if data.Timestamp != nil {
		capt.Timestamp = *data.Timestamp.Time
	}
	if data.Location != nil {
		capt.Location = &domain.Point{
			LAT:       data.Location.LAT,
			LNG:       data.Location.LNG,
			Elevation: data.Location.Elevation,
		}
	}
	if len(data.Tags) > 0 {
		capt.Tags = data.Tags
	}
}
