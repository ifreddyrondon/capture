package domain

import (
	"time"

	"github.com/ifreddyrondon/capture/pkg"
)

// User represents a user account.
type User struct {
	ID           string
	Email        string
	Password     []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Repositories []pkg.Repository
}
