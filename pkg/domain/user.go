package domain

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

// User represents a user account.
type User struct {
	ID           kallax.ULID `sql:"type:uuid" gorm:"primary_key"`
	Email        string      `sql:"not null" gorm:"unique_index"`
	Password     []byte
	CreatedAt    time.Time `sql:"not null"`
	UpdatedAt    time.Time `sql:"not null"`
	DeletedAt    *time.Time
	Repositories []Repository `gorm:"ForeignKey:UserID"`
}
