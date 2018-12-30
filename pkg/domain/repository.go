package domain

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

// Repository represent a place with the history of all captures.
type Repository struct {
	ID            kallax.ULID `json:"id" sql:"type:uuid,pk"`
	Name          string      `json:"name" sql:",notnull"`
	CurrentBranch string      `json:"current_branch" sql:",notnull"`
	Visibility    Visibility  `json:"visibility" sql:",notnull"`
	CreatedAt     time.Time   `json:"createdAt" sql:",notnull"`
	UpdatedAt     time.Time   `json:"updatedAt" sql:",notnull"`
	DeletedAt     *time.Time  `json:"-" pg:",soft_delete"`
	UserID        kallax.ULID `json:"owner" sql:"type:uuid"`
}
