package repository

// Service is the interface to be implemented by repository services
// It's the layer between HTTP server and Stores.
type Service interface {
	// Save a repository.
	Save(*Repository) error
}

// StoreService implementation of Service for repository
type StoreService struct {
	store Store
}

// NewService creates a new repository Service
func NewService(store Store) *StoreService {
	return &StoreService{store: store}
}

// Save a repository.
func (s *StoreService) Save(r *Repository) error {
	r.fillIfEmpty()
	return s.store.Save(r)
}
