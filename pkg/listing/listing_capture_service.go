package listing

import (
	"fmt"

	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

type notAuthorizedErr string

func (i notAuthorizedErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i notAuthorizedErr) IsNotAuthorized() bool { return true }

// CaptureStore provides access to the captures storage.
type CaptureStore interface {
	// List retrieve captures with domain.Listing attrs.
	List(*domain.Listing) ([]domain.Capture, error)
}

// CaptureService provides capture repository operations.
type CaptureService interface {
	// ListRepoCaptures list repo captures.
	ListRepoCaptures(*domain.User, *domain.Repository, *listing.Listing) (*ListCaptureResponse, error)
}

type captureService struct {
	s CaptureStore
}

// NewCaptureService creates a listing service with the necessary dependencies
func NewCaptureService(s CaptureStore) CaptureService {
	return &captureService{s: s}
}

func (s *captureService) ListRepoCaptures(u *domain.User, r *domain.Repository, l *listing.Listing) (*ListCaptureResponse, error) {
	if r.Visibility != domain.Public && r.UserID != u.ID {
		errStr := fmt.Sprintf("user %v not authorized to get captures from repo %v", u.ID, r.ID)
		return nil, notAuthorizedErr(errStr)
	}
	lcapt := domain.NewListing(*l)
	lcapt.Owner = &r.ID
	captures, err := s.s.List(lcapt)
	if err != nil {
		return nil, errors.Wrap(err, "err getting repo captures")
	}
	return newListCaptureResponse(captures, l), err
}

type ListCaptureResponse struct {
	Results []domain.Capture `json:"results"`
	Listing *listing.Listing `json:"listing"`
}

func newListCaptureResponse(repos []domain.Capture, l *listing.Listing) *ListCaptureResponse {
	if repos == nil {
		repos = make([]domain.Capture, 0)
	}
	return &ListCaptureResponse{Results: repos, Listing: l}
}
