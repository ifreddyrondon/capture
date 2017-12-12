package capture

import (
	"github.com/ifreddyrondon/gocapture/geocoding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const collection = "captures"

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Payload interface{}   `json:"payload"`
	*geocoding.Point
	*Timestamp
}

// NewCapture returns a new pointer to a Capture composed of the passed Point, Time and payload
func NewCapture(point *geocoding.Point, timestamp *Timestamp, payload interface{}) *Capture {
	return &Capture{
		ID:        bson.NewObjectId(),
		Payload:   payload,
		Point:     point,
		Timestamp: timestamp,
	}
}

// UnmarshalJSON decodes the capture from a JSON body.
// Throws an error if the body of the date cannot be interpreted by the JSON body.
// Implements the json.Unmarshaler Interface
func (c *Capture) UnmarshalJSON(data []byte) error {
	p := new(geocoding.Point)
	if err := p.UnmarshalJSON(data); err != nil {
		return err
	}

	t := Timestamp{}
	t.UnmarshalJSON(data)

	capture := NewCapture(p, &t, "")
	*c = *capture
	return nil
}

func (c *Capture) save(DB *mgo.Database) error {
	return DB.C(collection).Insert(c)
}
