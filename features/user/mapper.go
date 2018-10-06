package user

import (
	"github.com/ifreddyrondon/capture/features/user/decoder"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/src-d/go-kallax.v1"
)

func FromPostUser(postUser decoder.PostUser) (*User, error) {
	usr := User{
		ID:    kallax.NewULID(),
		Email: *postUser.Email,
	}

	if postUser.Password != nil {
		hash, err := hashPassword(*postUser.Password)
		if err != nil {
			return nil, err
		}
		usr.Password = hash
	}

	return &usr, nil
}

func hashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), 10)
}
