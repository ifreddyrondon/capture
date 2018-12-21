package domain

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

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
