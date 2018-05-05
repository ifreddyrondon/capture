package auth

import (
	"encoding/json"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/markbates/validate"
)

const (
	errEmailRequired = "email must not be blank!"
	errInvalidEmail  = "invalid email"
)

var errInvalidPayload = errors.New("cannot unmarshal json into valid credentials")

// BasicAuthCrendential represent the user credentials for basic authentication.
type BasicAuthCrendential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// alias for custom unmarhal
type basicAuthCrendentialJSON BasicAuthCrendential

func (bac *basicAuthCrendentialJSON) IsValid(errors *validate.Errors) {
	if bac.Email == "" {
		errors.Add("email", errEmailRequired)
	} else if !govalidator.IsEmail(bac.Email) {
		errors.Add("email", errInvalidEmail)
	}
	if bac.Password == "" {
		errors.Add("password", "password must not be blank!")
	}
}

// UnmarshalJSON decodes the BasicAuthCrendential from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (bac *BasicAuthCrendential) UnmarshalJSON(data []byte) error {
	var model basicAuthCrendentialJSON
	if err := json.Unmarshal(data, &model); err != nil {
		return errInvalidPayload
	}
	if err := validate.Validate(&model); err.HasAny() {
		return err
	}

	*bac = BasicAuthCrendential(model)
	return nil
}
