package capture

import (
	"encoding/json"
	"errors"
	"time"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/ifreddyrondon/gocapture/timestamp"
	"github.com/lib/pq"
	"github.com/mailru/easyjson/jwriter"

	"github.com/ifreddyrondon/gocapture/geocoding"
)

var (
	// ErrorBadPayload expected error when fails to unmarshal a capture
	ErrorBadPayload = errors.New("cannot unmarshal json into valid capture")
)

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID      kallax.ULID     `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Payload payload.Payload `json:"payload" sql:"not null;type:jsonb"`
	geocoding.Point
	Tags      pq.StringArray `json:"tags" sql:"not null;type:varchar(64)[]"`
	Timestamp time.Time      `json:"timestamp" sql:"not null"`
	CreatedAt time.Time      `json:"createdAt" sql:"not null"`
	UpdatedAt time.Time      `json:"updatedAt" sql:"not null"`
	DeletedAt *time.Time     `json:"-"`
}

// UnmarshalJSON decodes the capture from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (c *Capture) UnmarshalJSON(data []byte) error {
	var p geocoding.Point
	if err := p.UnmarshalJSON(data); err != nil {
		if err == geocoding.ErrUnmarshalPoint {
			return ErrorBadPayload
		}
		return err
	}
	c.Point = p

	var t timestamp.Timestamp
	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}
	c.Timestamp = t.Timestamp

	var payl payload.Payload
	if err := json.Unmarshal(data, &payl); err != nil {
		return err
	}
	c.Payload = payl

	var tags tagsJSON
	if err := tags.unmarshalJSON(data); err != nil {
		return err
	}
	c.Tags = tags.Tags

	return nil
}

// MarshalJSON supports json.Marshaler interface
func (c Capture) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCbca9c40EncodeGithubComIfreddyrondonGocaptureCapture(&w, c)
	return w.Buffer.BuildBytes(), w.Error
}
