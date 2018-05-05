package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/markbates/validate"

	"golang.org/x/crypto/bcrypt"

	kallax "gopkg.in/src-d/go-kallax.v1"
)

const (
	errEmailRequired = "email must not be blank!"
	errInvalidEmail  = "invalid email"
)

var errInvalidPayload = errors.New("cannot unmarshal json into valid user")

type emailDuplicateError struct {
	Email string
}

func (e *emailDuplicateError) Error() string {
	return fmt.Sprintf("email '%s' already exists", e.Email)
}

// User represents a user account.
type User struct {
	ID        kallax.ULID `json:"id" sql:"type:uuid" gorm:"primary_key"`
	Email     string      `json:"email" sql:"not null" gorm:"unique_index"`
	password  []byte
	CreatedAt time.Time  `json:"createdAt" sql:"not null"`
	UpdatedAt time.Time  `json:"updatedAt" sql:"not null"`
	DeletedAt *time.Time `json:"-"`
}

type userAlias User

type userJSON struct {
	userAlias
	Password string `json:"password"`
}

func (bac *userJSON) IsValid(errors *validate.Errors) {
	if bac.Email == "" {
		errors.Add("email", errEmailRequired)
	} else if !govalidator.IsEmail(bac.Email) {
		errors.Add("email", errInvalidEmail)
	}
}

// UnmarshalJSON decodes the user from a JSON body.
// Throws an error if the body cannot be interpreted.
// Implements the json.Unmarshaler Interface
func (u *User) UnmarshalJSON(data []byte) error {
	var user userJSON
	if err := json.Unmarshal(data, &user); err != nil {
		return errInvalidPayload
	}

	if err := validate.Validate(&user); err.HasAny() {
		return err
	}

	*u = User(user.userAlias)
	if len(user.Password) > 0 {
		if err := u.SetPassword([]byte(user.Password)); err != nil {
			return err
		}
	}
	return nil
}

// SetPassword stores a hashed version of a plain text password into the user.
func (u *User) SetPassword(pass []byte) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return err
	}
	u.password = hash
	return nil
}

// CheckPassword compares a hashed password with its possible plaintext equivalent.
func (u *User) CheckPassword(pass []byte) bool {
	if err := bcrypt.CompareHashAndPassword(u.password, pass); err != nil {
		return false
	}
	return true
}
