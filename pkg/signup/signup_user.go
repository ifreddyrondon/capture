package signup

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

// User represents a user response when sign-up an account.
type User struct {
	ID        kallax.ULID `json:"id"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"createdAt" sql:"not null"`
	UpdatedAt time.Time   `json:"updatedAt" sql:"not null"`
}
