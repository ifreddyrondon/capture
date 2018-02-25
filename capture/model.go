package capture

import (
	"time"

	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/ifreddyrondon/gocapture/timestamp"
	"github.com/mailru/easyjson/jwriter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Captures []Capture

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID      bson.ObjectId               `json:"id" bson:"_id,omitempty"`
	Payload *payload.ArrayNumberPayload `json:"payload"`
	*geocoding.Point
	*timestamp.Timestamp
	CreatedDate  time.Time `json:"created_date" bson:"createdDate"`
	LastModified time.Time `json:"last_modified" bson:"lastModified"`
}

// NewCapture returns a new pointer to a Capture composed of the passed Point, Time and payload
func NewCapture(point *geocoding.Point, timestamp *timestamp.Timestamp, payload *payload.ArrayNumberPayload) *Capture {
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
	var p geocoding.Point
	if err := p.UnmarshalJSON(data); err != nil {
		return err
	}

	var t timestamp.Timestamp
	t.UnmarshalJSON(data) // ignore err because timestamp always has a fallback

	// TODO: the payload could be of other types not only array number
	var payloadData payload.ArrayNumberPayload
	if err := payloadData.UnmarshalJSON(data); err != nil {
		return err
	}

	capture := NewCapture(&p, &t, &payloadData)
	*c = *capture
	return nil
}

// MarshalJSON supports json.Marshaler interface
func (c Capture) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComIfreddyrondonGocaptureCapture(&w, c)
	return w.Buffer.BuildBytes(), w.Error
}

func (c *Capture) create(DB *mgo.Database) error {
	now := time.Now()
	c.CreatedDate, c.LastModified = now, now
	return DB.C(Domain).Insert(c)
}

func (c *Capture) list(DB *mgo.Database, start, count int) (Captures, error) {
	results := Captures{}
	if err := DB.C(Domain).Find(nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}
