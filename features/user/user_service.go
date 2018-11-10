package user

import (
	"errors"
	"fmt"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/postgres"
	"gopkg.in/src-d/go-kallax.v1"
)

const uniqueConstraintEmail = "uix_users_email"

// ErrNotFound expected error when user is missing
var ErrNotFound = errors.New("user not found")

type emailDuplicateError struct {
	Email string
}

func (e *emailDuplicateError) Error() string {
	return fmt.Sprintf("email '%s' already exists", e.Email)
}

// GetterService get users
type GetterService interface {
	// Get a user by email
	GetByEmail(string) (*features.User, error)
	// Get a user by id
	GetByID(kallax.ULID) (*features.User, error)
}

// Service is the interface to be implemented by user services
// It's the layer between HTTP server and Stores.
type Service interface {
	// Save a collection.
	Save(*features.User) error
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
func (s *StoreService) Save(user *features.User) error {
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
func (s *StoreService) GetByEmail(email string) (*features.User, error) {
	return s.store.Get(StoreFilter{Email: &email})
}

// GetByID will look for a user with the same ID, or return
// user.ErrNotFound if no user is found.
func (s *StoreService) GetByID(id kallax.ULID) (*features.User, error) {
	return s.store.Get(StoreFilter{ID: &id})
}

type MockService struct {
	User *features.User
	Err  error
}

func (m *MockService) Save(user *features.User) error                  { return m.Err }
func (m *MockService) GetByEmail(email string) (*features.User, error) { return m.User, m.Err }
func (m *MockService) GetByID(id kallax.ULID) (*features.User, error)  { return m.User, m.Err }
