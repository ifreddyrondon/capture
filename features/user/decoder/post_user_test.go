package decoder_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/features/user/decoder"
	"github.com/stretchr/testify/assert"
)

func TestDecodePostUserOK(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name        string
		emailResult string
		body        string
	}{
		{
			name:        "decode user without password",
			emailResult: "test@localhost.com",
			body:        `{"email": "test@localhost.com"}`,
		},
		{
			name:        "decode user with password",
			emailResult: "test@localhost.com",
			body:        `{"email":"test@localhost.com","password":"1234"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var u decoder.PostUser
			err := decoder.Decode(r, &u)
			assert.Nil(t, err)
			assert.Equal(t, tc.emailResult, *u.Email)
		})
	}
}

func TestDecodePostUserError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "decode user's email missing",
			body: `{}`,
			err:  "email must not be blank",
		},
		{
			name: "decode user's email missing",
			body: `{"email": "test@"}`,
			err:  "invalid email",
		},
		{
			name: "decode user's password too short",
			body: `{"email":"test@localhost.com","password":"1"}`,
			err:  "password must have at least four characters",
		},
		{
			name: "invalid user payload",
			body: `.`,
			err:  "cannot unmarshal json into valid user",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var u decoder.PostUser
			err := decoder.Decode(r, &u)
			assert.EqualError(t, err, tc.err)
		})
	}
}
