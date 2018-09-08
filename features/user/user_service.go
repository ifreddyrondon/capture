package user

import (
	"github.com/ifreddyrondon/capture/features/postgres"
	"gopkg.in/src-d/go-kallax.v1"
)

const uniqueConstraintEmail = "uix_users_email"

// GetterService get users
type GetterService interface {
	// Get a user by email
	GetByEmail(string) (*User, error)
	// Get a user by id
	GetByID(kallax.ULID) (*User, error)
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

// GetByEmail will look for a user with the same email address, or return
// user.ErrNotFound if no user is found.
//
// ByEmail is NOT case sensitive.
func (s *StoreService) GetByEmail(email string) (*User, error) {
	return s.store.GetByEmail(email)
}

// GetByID will look for a user with the same ID, or return
// user.ErrNotFound if no user is found.
func (s *StoreService) GetByID(id kallax.ULID) (*User, error) {
	return s.store.GetByID(id)
}
