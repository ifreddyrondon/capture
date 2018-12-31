package domain

import (
	"time"

	"gopkg.in/src-d/go-kallax.v1"
)

// Metric represent a captured value from a device/sensor
type Metric struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Payload represent an association of metrics
type Payload []Metric

// Point represents a physical Point in geographic notation [lat, lng].
type Point struct {
	LAT       *float64 `json:"lat"`
	LNG       *float64 `json:"lng"`
	Elevation *float64 `json:"elevation,omitempty"`
}

type Capture struct {
	ID           kallax.ULID `json:"id" sql:"type:uuid,pk"`
	Payload      Payload     `json:"payload" sql:"type:jsonb,notnull"`
	Location     *Point      `json:"location" sql:"type:jsonb"`
	Tags         []string    `json:"tags" sql:",array,notnull"`
	Timestamp    time.Time   `json:"timestamp" sql:",notnull"`
	CreatedAt    time.Time   `json:"createdAt" sql:",notnull"`
	UpdatedAt    time.Time   `json:"updatedAt" sql:",notnull"`
	DeletedAt    *time.Time  `json:"-" pg:",soft_delete"`
	RepositoryID kallax.ULID `json:"repoId" sql:"type:uuid"`
}
