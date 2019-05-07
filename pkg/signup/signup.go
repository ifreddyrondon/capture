package signup

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/validate"
)

const (
	errEmailRequired      = "email must not be blank"
	errInvalidEmail       = "invalid email"
	errInvalidPasswordLen = "password must have at least four characters"

	minPasswordLen = 4
)

type Payload struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u Payload) Validate() error {
	e := validate.NewErrors()
	if u.Email == nil {
		e.Add("email", errEmailRequired)
	} else if !govalidator.IsEmail(*u.Email) {
		e.Add("email", errInvalidEmail)
	}

	if u.Password != nil {
		if len(*u.Password) < minPasswordLen {
			e.Add("password", errInvalidPasswordLen)
		}
	}

	if e.HasAny() {
		return e
	}

	return nil
}
