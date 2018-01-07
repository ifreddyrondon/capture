package app

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/timestamp"
	"gopkg.in/mgo.v2/bson"
)

const (
	BranchDomain  = "branches"
	CaptureDomain = "captures"
)

type Router interface {
	Pattern() string
	Router() chi.Router
}

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Payload interface{}   `json:"payload"`
	*geocoding.Point
	*timestamp.Timestamp
	CreatedDate  time.Time `json:"created_date" bson:"createdDate"`
	LastModified time.Time `json:"last_modified" bson:"lastModified"`
}

type CaptureService interface {
	CreateCapture(*Capture) error
}
