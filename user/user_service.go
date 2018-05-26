package user

import "github.com/ifreddyrondon/gocapture/postgres"

const uniqueConstraintEmail = "uix_users_email"

// GetterService get users
type GetterService interface {
	// Get a user by email
	Get(string) (*User, error)
}

// Service is the interface to be implemented by user services
// It's the layer between HTTP server and Repositories.
type Service interface {
	// Save a collection.
	Save(*User) error
	// Get a user by email
	GetterService
}

// REPOService implementation of Service for repository
type REPOService struct {
	repository Repository
}

// NewService creates a new user Service
func NewService(repository Repository) *REPOService {
	return &REPOService{repository: repository}
}

// Save a capture
func (s *REPOService) Save(user *User) error {
	user.fillIfEmpty()
	err := s.repository.Save(user)

	if err != nil {
		if postgres.IsUniqueConstraintError(err, uniqueConstraintEmail) {
			return &emailDuplicateError{Email: user.Email}
		}
		return err
	}
	return nil
}

// Get a user by email
func (s *REPOService) Get(email string) (*User, error) {
	return s.repository.Get(email)
}
