package features

import (
	"time"

	"github.com/ifreddyrondon/capture/features/capture/geocoding"
	"github.com/ifreddyrondon/capture/features/capture/payload"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/src-d/go-kallax.v1"
)

var visibilityTypes = [...]string{"public", "private"}

type Visibility string

func AllowedVisibility(test string) bool {
	if test == "" {
		return false
	}

	for i := range visibilityTypes {
		if visibilityTypes[i] == test {
			return true
		}
	}

	return false
}

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

// Branch is a partial or full collection of captures within a repository.
type Branch struct {
	ID       kallax.ULID `json:"id"`
	Name     string      `json:"name"`
	Captures []Capture   `json:"captures"`
}

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID        kallax.ULID      `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Payload   payload.Payload  `json:"payload" sql:"not null;type:jsonb"`
	Location  *geocoding.Point `json:"location" sql:"type:jsonb"`
	Tags      pq.StringArray   `json:"tags" sql:"not null;type:varchar(64)[]"`
	Timestamp time.Time        `json:"timestamp" sql:"not null"`
	CreatedAt time.Time        `json:"createdAt" sql:"not null"`
	UpdatedAt time.Time        `json:"updatedAt" sql:"not null"`
	DeletedAt *time.Time       `json:"-"`
}

// Repository represent a place with the history of all captures.
type Repository struct {
	ID            kallax.ULID `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Name          string      `json:"name"`
	CurrentBranch string      `json:"current_branch"`
	Visibility    Visibility  `json:"visibility"`
	CreatedAt     time.Time   `json:"createdAt" sql:"not null"`
	UpdatedAt     time.Time   `json:"updatedAt" sql:"not null"`
	DeletedAt     *time.Time  `json:"-"`
	UserID        kallax.ULID `json:"owner"`
}

// User represents a user account.
type User struct {
	ID           kallax.ULID  `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Email        string       `json:"email" sql:"not null" gorm:"unique_index"`
	Password     []byte       `json:"-"`
	CreatedAt    time.Time    `json:"createdAt" sql:"not null"`
	UpdatedAt    time.Time    `json:"updatedAt" sql:"not null"`
	DeletedAt    *time.Time   `json:"-"`
	Repositories []Repository `json:"-" gorm:"ForeignKey:UserID"`
}

// CheckPassword compares a hashed password with its possible plaintext equivalent.
func (u *User) CheckPassword(pass string) bool {
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pass)); err != nil {
		return false
	}
	return true
}
