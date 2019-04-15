package authenticating

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/validate"

	"github.com/ifreddyrondon/capture/pkg/validator"
)

const (
	errEmailRequired    = "email must not be blank"
	errInvalidEmail     = "invalid email"
	errPasswordRequired = "password must not be blank"
)

// Validator for authenticating request
const Validator validator.StringValidator = "cannot unmarshal body into valid user credentials"

// BasicCredential for authentication purposes.
type BasicCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// OK implementation of validator.OK
func (c *BasicCredential) OK() error {
	e := validate.NewErrors()

	if c.Email == "" {
		e.Add("email", errEmailRequired)
	} else if !govalidator.IsEmail(c.Email) {
		e.Add("email", errInvalidEmail)
	}
	if c.Password == "" {
		e.Add("password", errPasswordRequired)
	}

	if e.HasAny() {
		return e
	}
	return nil
}
