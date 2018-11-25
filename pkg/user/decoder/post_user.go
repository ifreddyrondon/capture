package decoder

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/markbates/validate"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/src-d/go-kallax.v1"
)

const (
	errEmailRequired      = "email must not be blank"
	errInvalidEmail       = "invalid email"
	errInvalidPasswordLen = "password must have at least four characters"

	minPasswordLen = 4
)

type PostUser struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u PostUser) OK() error {
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

func (u PostUser) User(usr *pkg.User) error {
	usr.ID = kallax.NewULID()
	usr.Email = *u.Email
	now := time.Now()
	usr.CreatedAt = now
	usr.UpdatedAt = now

	if u.Password != nil {
		hash, err := hashPassword(*u.Password)
		if err != nil {
			return err
		}
		usr.Password = hash
	}

	return nil
}

func hashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), 10)
}
