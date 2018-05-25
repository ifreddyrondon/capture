package collection

import (
	"errors"
	"time"

	"github.com/markbates/validate"

	kallax "gopkg.in/src-d/go-kallax.v1"
)

var errInvalidPayload = errors.New("cannot unmarshal json into valid collection")

// Collection represent a collection of captures.
type Collection struct {
	ID            kallax.ULID `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Name          string      `json:"name"`
	CurrentBranch string      `json:"current_branch"`
	CreatedAt     time.Time   `json:"createdAt" sql:"not null"`
	UpdatedAt     time.Time   `json:"updatedAt" sql:"not null"`
	DeletedAt     *time.Time  `json:"-"`
}

// UnmarshalJSON decodes the user from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (c *Collection) UnmarshalJSON(data []byte) error {
	var model jsonCollection
	if err := model.UnmarshalJSON(data); err != nil {
		return errInvalidPayload
	}
	if err := validate.Validate(&model); err.HasAny() {
		return err
	}

	c.Name = model.Name
	return nil
}

func (c *Collection) fillIfEmpty() {
	if c.ID.IsEmpty() {
		c.ID = kallax.NewULID()
	}
}
