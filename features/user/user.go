package user

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/src-d/go-kallax.v1"
)

// ErrNotFound expected error when user is missing
var ErrNotFound = errors.New("user not found")

type emailDuplicateError struct {
	Email string
}

func (e *emailDuplicateError) Error() string {
	return fmt.Sprintf("email '%s' already exists", e.Email)
}

// User represents a user account.
type User struct {
	ID        kallax.ULID `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Email     string      `json:"email" sql:"not null" gorm:"unique_index"`
	Password  []byte      `json:"-"`
	CreatedAt time.Time   `json:"createdAt" sql:"not null"`
	UpdatedAt time.Time   `json:"updatedAt" sql:"not null"`
	DeletedAt *time.Time  `json:"-"`
}

// CheckPassword compares a hashed password with its possible plaintext equivalent.
func (u *User) CheckPassword(pass string) bool {
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pass)); err != nil {
		return false
	}
	return true
}
