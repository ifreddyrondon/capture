package features

import (
	"time"

	"github.com/ifreddyrondon/capture/features/capture"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/src-d/go-kallax.v1"
)

// Repository represent a place with the history of all captures.
type Repository struct {
	ID            kallax.ULID `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Name          string      `json:"name"`
	CurrentBranch string      `json:"current_branch"`
	Shared        bool        `json:"shared"`
	CreatedAt     time.Time   `json:"createdAt" sql:"not null"`
	UpdatedAt     time.Time   `json:"updatedAt" sql:"not null"`
	DeletedAt     *time.Time  `json:"-"`
	UserID        kallax.ULID `json:"owner"`
}

// Branch is a partial or full collection of captures within a repository.
type Branch struct {
	ID       kallax.ULID       `json:"id"`
	Name     string            `json:"name"`
	Captures []capture.Capture `json:"captures"`
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
