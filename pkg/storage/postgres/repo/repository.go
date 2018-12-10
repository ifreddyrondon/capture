package repo

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

type Repository struct {
	ID            kallax.ULID `sql:"type:uuid" gorm:"primary_key"`
	Name          string
	CurrentBranch string
	Visibility    string
	CreatedAt     time.Time `sql:"not null"`
	UpdatedAt     time.Time `sql:"not null"`
	DeletedAt     *time.Time
	UserID        kallax.ULID
}
