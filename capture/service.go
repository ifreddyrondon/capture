package capture

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/payload/numberlist"
	"gopkg.in/mgo.v2"
)

// Service is the interface implemented by capture
// that can make CRUD operations over captures.
type Service interface {
	// Create a capture into the database.
	Create(*geocoding.Point, time.Time, *numberlist.Payload) (*Capture, error)
	// List retrieve the count captures from start index.
	List(start, count int) (Captures, error)
	// Get retrive a capture by id
	Get(string) (*Capture, error)
}

// MgoService implementation of capture.Service for
// Mongo database.
type MgoService struct {
	*mgo.Collection
}

// Create a capture into the database.
func (s *MgoService) Create(p *geocoding.Point, t time.Time, payload *numberlist.Payload) (*Capture, error) {
	c := New(p, t, payload)
	now := time.Now()
	c.CreatedDate, c.LastModified = now, now
	if err := s.Collection.Insert(c); err != nil {
		return nil, err
	}
	return c, nil
}

// List retrieve the count captures from start index.
func (s *MgoService) List(start, count int) (Captures, error) {
	results := Captures{}
	if err := s.Collection.Find(nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// Get retrive a capture by id
func (s *MgoService) Get(id string) (*Capture, error) {
	var result Capture
	if err := s.Collection.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
