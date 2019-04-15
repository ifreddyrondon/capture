package creating

import (
	"strings"
	"time"

	"github.com/gobuffalo/validate"

	"github.com/ifreddyrondon/capture/pkg/validator"
)

const (
	errNameRequired         = "name must not be blank"
	errVisibilityNotAllowed = "not allowed visibility type. it Could be one of public, or private. Default: public"
)

// Validator for sign-up request payload
const Validator validator.StringValidator = "cannot unmarshal json into valid repository"

var visibilityTypes = [...]string{"public", "private"}

type Visibility string

func isAllowedVisibility(test string) bool {
	if test == "" {
		return false
	}
	for i := range visibilityTypes {
		if visibilityTypes[i] == test {
			return true
		}
	}
	return false
}

var (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type Payload struct {
	Name       *string `json:"name"`
	Visibility *string `json:"visibility"`
}

func (p Payload) OK() error {
	e := validate.NewErrors()
	if p.Name == nil {
		e.Add("name", errNameRequired)
	} else if len(strings.TrimSpace(*p.Name)) == 0 {
		e.Add("name", errNameRequired)
	}

	if p.Visibility != nil {
		if !isAllowedVisibility(*p.Visibility) {
			e.Add("name", errVisibilityNotAllowed)
		}
	}
	if e.HasAny() {
		return e
	}

	return nil
}

type Repository struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"createdAt" `
	UpdatedAt  time.Time `json:"updatedAt" `
}
