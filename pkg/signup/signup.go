package signup

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/validate"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

const (
	errEmailRequired      = "email must not be blank"
	errInvalidEmail       = "invalid email"
	errInvalidPasswordLen = "password must have at least four characters"

	minPasswordLen = 4
)

// Validator for sign-up request payload
const Validator validator.StringValidator = "cannot unmarshal json into valid user"

type Payload struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u Payload) OK() error {
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
