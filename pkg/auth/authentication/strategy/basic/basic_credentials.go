package basic

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/markbates/validate"
)

const (
	errEmailRequired    = "email must not be blank"
	errInvalidEmail     = "invalid email"
	errPasswordRequired = "password must not be blank"
)

var errInvalidPayload = errors.New("cannot unmarshal json into valid credentials")

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Credentials represent the user credentials for basic authentication.
type Credentials credentials

// IsValid validates if credentials are valid.
// Implements Validator from github.com/markbates/validate
func (c *Credentials) IsValid(errors *validate.Errors) {
	if c.Email == "" {
		errors.Add("email", errEmailRequired)
	} else if !govalidator.IsEmail(c.Email) {
		errors.Add("email", errInvalidEmail)
	}
	if c.Password == "" {
		errors.Add("password", errPasswordRequired)
	}
}

// UnmarshalJSON decodes the BasicAuthCredentials from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (c *Credentials) UnmarshalJSON(data []byte) error {
	var model credentials
	if err := model.UnmarshalJSON(data); err != nil {
		return errInvalidPayload
	}
	*c = Credentials(model)
	if err := validate.Validate(c); err.HasAny() {
		return err
	}

	return nil
}
