package decoder_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/user/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestDecodeFromPostUserOK(t *testing.T) {
	email, pass := "test@example.com", "1234"

	t.Parallel()
	tt := []struct {
		name     string
		body     string
		expected decoder.PostUser
	}{
		{
			name:     "decode user without password",
			body:     `{"email": "test@example.com"}`,
			expected: decoder.PostUser{Email: &email},
		},
		{
			name:     "decode user with password",
			body:     `{"email":"test@example.com","password":"1234"}`,
			expected: decoder.PostUser{Email: &email, Password: &pass},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var u decoder.PostUser
			err := decoder.Decode(r, &u)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Email, u.Email)
			assert.Equal(t, tc.expected.Password, u.Password)
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
			name: "decode user's invalid missing",
			body: `{"email": "test@"}`,
			err:  "invalid email",
		},
		{
			name: "decode user's password too short",
			body: `{"email":"test@example.com","password":"1"}`,
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

func TestUserPostUserOK(t *testing.T) {
	email, pass := "test@example.com", "1234"
	t.Parallel()
	tt := []struct {
		name     string
		postUser decoder.PostUser
		expected pkg.User
	}{
		{
			name:     "get user from postUser without password",
			postUser: decoder.PostUser{Email: &email},
			expected: pkg.User{Email: email},
		},
		{
			name:     "get user from postUser with password",
			postUser: decoder.PostUser{Email: &email, Password: &pass},
			expected: pkg.User{Email: email},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var u pkg.User
			err := tc.postUser.User(&u)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Email, u.Email)
			// test user fields filled with not default values
			assert.NotEqual(t, kallax.ULID{}, u.ID)
			assert.NotEqual(t, time.Time{}, u.CreatedAt)
			assert.NotEqual(t, time.Time{}, u.UpdatedAt)
		})
	}
}
