package pkg

import (
	"time"

	"github.com/ifreddyrondon/capture/pkg/capture/geocoding"
	"github.com/ifreddyrondon/capture/pkg/capture/payload"
	"github.com/lib/pq"
	"gopkg.in/src-d/go-kallax.v1"
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
