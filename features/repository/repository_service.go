package repository

import "github.com/ifreddyrondon/capture/features"

// Service is the interface to be implemented by repository services
// It's the layer between HTTP server and Stores.
type Service interface {
	// Save a repository.
	Save(user *features.Repository) error
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
func (s *StoreService) Save(r *features.Repository) error {
	return s.store.Save(r)
}
