package capture

import (
	"github.com/ifreddyrondon/capture/features"
	kallax "gopkg.in/src-d/go-kallax.v1"
)

// Service is the interface to be implemented by capture services
// It's the layer between HTTP server and Stores.
type Service interface {
	// Save a capture.
	Save(*features.Capture) error
	// SaveBulk captures.
	SaveBulk(...*features.Capture) (Captures, error)
	// List retrieve captures from start index to count.
	List(start, count int) (Captures, error)
	// Get a capture by id
	Get(kallax.ULID) (*features.Capture, error)
	// Delete a capture by id
	Delete(*features.Capture) error
	// Update a capture from an updated one, will only update those changed & non blank fields.
	Update(original *features.Capture, updates features.Capture) error
}

// StoreService implementation of Service for capture
type StoreService struct {
	store Store
}

// NewService creates a new user Service
func NewService(store Store) *StoreService {
	return &StoreService{store: store}
}

// Save a capture.
func (s *StoreService) Save(capt *features.Capture) error {
	return s.store.Save(capt)
}

// SaveBulk captures.
func (s *StoreService) SaveBulk(captures ...*features.Capture) (Captures, error) {
	return s.store.SaveBulk(captures...)
}

// List retrieve the count captures from start index.
func (s *StoreService) List(start, count int) (Captures, error) {
	return s.store.List(start, count)
}

// Get a capture by id
func (s *StoreService) Get(id kallax.ULID) (*features.Capture, error) {
	return s.store.Get(id)
}

// Delete a capture by id
func (s *StoreService) Delete(capt *features.Capture) error {
	return s.store.Delete(capt)
}

// Update a capture
func (s *StoreService) Update(original *features.Capture, updates features.Capture) error {
	return s.store.Update(original, updates)
}
