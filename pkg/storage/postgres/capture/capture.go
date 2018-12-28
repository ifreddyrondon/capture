package capture

import (
	"time"

	"github.com/lib/pq"
	"gopkg.in/src-d/go-kallax.v1"
)

type Capture struct {
	ID        kallax.ULID    `sql:"type:uuid" gorm:"primary_key"`
	Payload   payload        `sql:"not null;type:jsonb"`
	Location  *point         `sql:"type:jsonb"`
	Tags      pq.StringArray `sql:"not null" gorm:"type:varchar(64)[]"`
	Timestamp time.Time      `sql:"not null"`
	CreatedAt time.Time      `sql:"not null"`
	UpdatedAt time.Time      `sql:"not null"`
	DeletedAt *time.Time
}
