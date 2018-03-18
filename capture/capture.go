package capture

import (
	"time"

	"github.com/ifreddyrondon/gocapture/payload"

	"github.com/ifreddyrondon/gocapture/geocoding"
	"github.com/ifreddyrondon/gocapture/timestamp"
	"github.com/mailru/easyjson/jwriter"
	"gopkg.in/mgo.v2/bson"
)

type Captures []Capture

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
	*numberlist.Payload `json:"payload"`
	*geocoding.Point
	Timestamp    time.Time
	CreatedDate  time.Time `json:"created_date" bson:"createdDate"`
	LastModified time.Time `json:"last_modified" bson:"lastModified"`
}

// New returns a new pointer to a Capture composed of the passed Point, Time and payload
func New(point *geocoding.Point, timestamp time.Time, payload *numberlist.Payload) *Capture {
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
	var payload numberlist.Payload
	if err := payload.UnmarshalJSON(data); err != nil {
		return err
	}

	capture := New(&p, t.Timestamp, &payload)
	*c = *capture
	return nil
}

// MarshalJSON supports json.Marshaler interface
func (c Capture) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComIfreddyrondonGocaptureCapture(&w, c)
	return w.Buffer.BuildBytes(), w.Error
}
