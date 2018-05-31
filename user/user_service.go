package user

import "github.com/ifreddyrondon/gocapture/postgres"

const uniqueConstraintEmail = "uix_users_email"

// GetterService get users
type GetterService interface {
	// Get a user by email
	Get(string) (*User, error)
}

// Service is the interface to be implemented by user services
// It's the layer between HTTP server and Stores.
type Service interface {
	// Save a collection.
	Save(*User) error
	// Get a user by email
	GetterService
}

// StoreService implementation of Service for user
type StoreService struct {
	store Store
}

// NewService creates a new user Service
func NewService(store Store) *StoreService {
	return &StoreService{store: store}
}

// Save a capture
func (s *StoreService) Save(user *User) error {
	user.fillIfEmpty()
	err := s.store.Save(user)

	if err != nil {
		if postgres.IsUniqueConstraintError(err, uniqueConstraintEmail) {
			return &emailDuplicateError{Email: user.Email}
		}
		return err
	}
	return nil
}

// Get a user by email
func (s *StoreService) Get(email string) (*User, error) {
	return s.store.Get(email)
}
