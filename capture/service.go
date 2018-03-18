package capture

import (
	"time"

	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/payload"
	"gopkg.in/mgo.v2"
)

// Service is the interface implemented by capture
// that can make CRUD operations over captures.
type Service interface {
	Create(*geocoding.Point, time.Time, *numberlist.Payload) (*Capture, error)
	List(start, count int) (Captures, error)
}

// MgoService implementation of capture.Service for
// Mongo database.
type MgoService struct {
	DB *mgo.Database
}

// Create generate a capture into the database.
func (s *MgoService) Create(p *geocoding.Point, t time.Time, payload *numberlist.Payload) (*Capture, error) {
	c := New(p, t, payload)
	now := time.Now()
	c.CreatedDate, c.LastModified = now, now
	if err := s.DB.C(Domain).Insert(c); err != nil {
		return nil, err
	}
	return c, nil
}

// List retrieve the count captures from start index.
func (s *MgoService) List(start, count int) (Captures, error) {
	results := Captures{}
	if err := s.DB.C(Domain).Find(nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}
