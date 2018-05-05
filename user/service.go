package user

import (
	"github.com/ifreddyrondon/gocapture/postgres"
	"github.com/jinzhu/gorm"
	kallax "gopkg.in/src-d/go-kallax.v1"
)

const uniqueConstraintEmail = "uix_users_email"

// Service is the interface implemented by user
// It make CRUD operations over users.
type Service interface {
	// Save user into the database.
	Save(*User) error
	// Get a user by email
	Get(string) (*User, error)
	// Delete a user
	// Delete(*User) error
	// Update a user from an updated one, will only update those changed & non blank fields.
	// Update(original *User, updates User) error
}

// PGService implementation of user.Service for Postgres database.
type PGService struct {
	DB *gorm.DB
}

// Migrate (panic) runs schema migration.
func (pgs *PGService) Migrate() {
	pgs.DB.AutoMigrate(User{})
}

// Drop (panic) delete schema.
func (pgs *PGService) Drop() {
	pgs.DB.DropTableIfExists(User{})
}

// Save capture into the database.
func (pgs *PGService) Save(user *User) error {
	user.ID = kallax.NewULID()
	err := pgs.DB.Create(user).Error

	if err != nil {
		if postgres.IsUniqueConstraintError(err, uniqueConstraintEmail) {
			return &emailDuplicateError{Email: user.Email}
		}
		return err
	}
	return nil
}

// Get a user by email
func (pgs *PGService) Get(email string) (*User, error) {
	var result User
	if pgs.DB.Where(&User{Email: email}).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}
