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
	ID        kallax.ULID `json:"id"`
	Payload   Payload     `json:"payload"`
	Location  *Point      `json:"location"`
	Tags      []string    `json:"tags"`
	Timestamp time.Time   `json:"timestamp"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt *time.Time  `json:"-"`
}
