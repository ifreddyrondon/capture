package domain

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

// User represents a user account.
type User struct {
	ID        kallax.ULID `sql:"type:uuid,pk"`
	Email     string      `sql:",notnull,unique"`
	Password  []byte      `sql:",notnull"`
	CreatedAt time.Time   `sql:",notnull"`
	UpdatedAt time.Time   `sql:",notnull"`
	DeletedAt *time.Time  `pg:",soft_delete"`
}
