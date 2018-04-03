package capture

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ifreddyrondon/gocapture/payload"
	"github.com/ifreddyrondon/gocapture/timestamp"
	jwriter "github.com/mailru/easyjson/jwriter"

	"github.com/ifreddyrondon/gocapture/geocoding"
)

var (
	// ErrorBadPayload expected error when fails to unmarshal a capture
	ErrorBadPayload = errors.New("cannot unmarshal json into valid capture")
)

// Capture is the representation of data sample of any kind taken at a specific time and location.
type Capture struct {
	ID      uint64          `json:"id" gorm:"primary_key"`
	Payload payload.Payload `json:"payload" sql:"type:jsonb"`
	geocoding.Point
	Timestamp time.Time  `json:"timestamp"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// New returns a new pointer to a Capture composed by a payload a timestamp and a point.
func New(payl payload.Payload, timestamp time.Time, point geocoding.Point) (*Capture, error) {
	if len(payl) == 0 {
		return nil, payload.ErrorMissingPayload
	}

	return &Capture{
		Payload:   payl,
		Point:     point,
		Timestamp: timestamp,
	}, nil
}

// UnmarshalJSON decodes the capture from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (c *Capture) UnmarshalJSON(data []byte) error {
	var p geocoding.Point
	if err := p.UnmarshalJSON(data); err != nil {
		if err == geocoding.ErrorUnmarshalPoint {
			return ErrorBadPayload
		}
		return err
	}

	var t timestamp.Timestamp
	t.UnmarshalJSON(data) // ignore err because timestamp always has a fallback

	var payl payload.Payload
	if err := json.Unmarshal(data, &payl); err != nil {
		return err
	}

	capture, err := New(payl, t.Timestamp, p)
	if err != nil {
		return err
	}
	*c = *capture
	return nil
}

// MarshalJSON supports json.Marshaler interface
func (c Capture) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCbca9c40EncodeGithubComIfreddyrondonGocaptureCapture(&w, c)
	return w.Buffer.BuildBytes(), w.Error
}
