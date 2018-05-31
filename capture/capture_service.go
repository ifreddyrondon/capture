package capture

import kallax "gopkg.in/src-d/go-kallax.v1"

// Service is the interface to be implemented by capture services
// It's the layer between HTTP server and Repositories.
type Service interface {
	// Save a capture.
	Save(*Capture) error
	// SaveBulk captures.
	SaveBulk(...*Capture) (Captures, error)
	// List retrieve captures from start index to count.
	List(start, count int) (Captures, error)
	// Get a capture by id
	Get(kallax.ULID) (*Capture, error)
	// Delete a capture by id
	Delete(*Capture) error
	// Update a capture from an updated one, will only update those changed & non blank fields.
	Update(original *Capture, updates Capture) error
}

// REPOService implementation of Service for repository
type REPOService struct {
	repository Repository
}

// NewService creates a new user Service
func NewService(repository Repository) *REPOService {
	return &REPOService{repository: repository}
}

// Save a capture.
func (s *REPOService) Save(capt *Capture) error {
	capt.fillIfEmpty()
	return s.repository.Save(capt)
}

// SaveBulk captures.
func (s *REPOService) SaveBulk(captures ...*Capture) (Captures, error) {
	return s.repository.SaveBulk(captures...)
}

// List retrieve the count captures from start index.
func (s *REPOService) List(start, count int) (Captures, error) {
	return s.repository.List(start, count)
}

// Get a capture by id
func (s *REPOService) Get(id kallax.ULID) (*Capture, error) {
	return s.repository.Get(id)
}

// Delete a capture by id
func (s *REPOService) Delete(capt *Capture) error {
	return s.repository.Delete(capt)
}

// Update a capture
func (s *REPOService) Update(original *Capture, updates Capture) error {
	return s.repository.Update(original, updates)
}
