package getting

import (
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

// CaptureStore provides access to the capture storage.
type CaptureStore interface {
	// Get retrieve a capture from storage.
	Get(captureID, repoID kallax.ULID) (*domain.Capture, error)
}

// CaptureService provides getting capture operations.
type CaptureService interface {
	// Get retrieve a repo capture if belong to an user or if it's a public repo.
	Get(kallax.ULID, *domain.Repository) (*domain.Capture, error)
}

type captureService struct {
	s CaptureStore
}

// NewCaptureService creates a getting service with the necessary dependencies
func NewCaptureService(s CaptureStore) CaptureService {
	return &captureService{s: s}
}

func (s *captureService) Get(id kallax.ULID, r *domain.Repository) (*domain.Capture, error) {
	capt, err := s.s.Get(id, r.ID)
	if err != nil {
		return nil, errors.Wrap(err, "could not get capture")
	}
	return capt, nil
}
