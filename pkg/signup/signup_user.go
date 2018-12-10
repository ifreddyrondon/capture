package signup

import (
	"time"
)

// User represents a user response when sign-up an account.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt" sql:"not null"`
	UpdatedAt time.Time `json:"updatedAt" sql:"not null"`
}
