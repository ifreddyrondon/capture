package basic

import (
	"encoding/json"
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

// Crendentials represent the user credentials for basic authentication.
type Crendentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// alias for custom unmarhal
type crendentialJSON Crendentials

func (c *crendentialJSON) IsValid(errors *validate.Errors) {
	if c.Email == "" {
		errors.Add("email", errEmailRequired)
	} else if !govalidator.IsEmail(c.Email) {
		errors.Add("email", errInvalidEmail)
	}
	if c.Password == "" {
		errors.Add("password", errPasswordRequired)
	}
}

// UnmarshalJSON decodes the BasicAuthCrendential from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (c *Crendentials) UnmarshalJSON(data []byte) error {
	var model crendentialJSON
	if err := json.Unmarshal(data, &model); err != nil {
		return errInvalidPayload
	}
	if err := validate.Validate(&model); err.HasAny() {
		return err
	}

	*c = Crendentials(model)
	return nil
}
