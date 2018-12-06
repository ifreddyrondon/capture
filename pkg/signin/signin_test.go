package signin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/signin"
	"github.com/stretchr/testify/assert"
)

func TestValidatePayloadOK(t *testing.T) {
	email, pass := "test@example.com", "1234"

	t.Parallel()
	tt := []struct {
		name     string
		body     string
		expected signin.Payload
	}{
		{
			name:     "decode user without password",
			body:     `{"email": "test@example.com"}`,
			expected: signin.Payload{Email: &email},
		},
		{
			name:     "decode user with password",
			body:     `{"email":"test@example.com","password":"1234"}`,
			expected: signin.Payload{Email: &email, Password: &pass},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p signin.Payload
			err := signin.Validator.Decode(r, &p)
			assert.Nil(t, err)
			assert.Equal(t, p.Email, p.Email)
			assert.Equal(t, p.Password, p.Password)
		})
	}
}

func TestValidatePayloadError(t *testing.T) {
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

			var p signin.Payload
			err := signin.Validator.Decode(r, &p)
			assert.EqualError(t, err, tc.err)
		})
	}
}
