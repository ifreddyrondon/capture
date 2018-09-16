package basic_test

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/features/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/features/user"
	"github.com/stretchr/testify/assert"
)

const (
	userEmail    = "test@example.com"
	userPassword = "secret"
	hashedPass   = "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK"
)

func TestValidateSuccess(t *testing.T) {
	t.Parallel()
	strategy := basic.New(&user.MockService{User: &user.User{Email: userEmail, Password: []byte(hashedPass)}})

	body := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, userEmail, userPassword))
	req := httptest.NewRequest("GET", "/", body)

	u, err := strategy.Validate(req)
	assert.Nil(t, err)
	assert.Equal(t, userEmail, u.Email)
}

func TestValidateInvalidCredentials(t *testing.T) {
	t.Parallel()
	strategy := basic.New(&user.MockService{User: &user.User{Email: userEmail, Password: []byte(hashedPass)}})

	tt := []struct {
		name string
		body io.Reader
		errs []string
	}{
		{
			name: "invalid credentials",
			body: strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, userEmail, "123")),
			errs: []string{"invalid email or password"},
		},
		{
			name: "missing email",
			body: strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, "bla@example.com", "123")),
			errs: []string{"invalid email or password"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", tc.body)
			_, err := strategy.Validate(req)
			assert.Error(t, err)
			assert.True(t, strategy.IsErrCredentials(err))
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

func TestValidateFailsDecoding(t *testing.T) {
	t.Parallel()
	strategy := basic.New(&user.MockService{Err: errors.New("test")})

	tt := []struct {
		name string
		body io.Reader
		errs []string
	}{
		{
			name: "invalid json",
			body: strings.NewReader("{"),
			errs: []string{"unexpected EOF"},
		},
		{
			name: "missing data",
			body: strings.NewReader("{}"),
			errs: []string{"email must not be blank", "password must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", tc.body)
			_, err := strategy.Validate(req)
			assert.Error(t, err)
			assert.True(t, strategy.IsErrDecoding(err))
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

func TestValidateFailsUnknownErr(t *testing.T) {
	t.Parallel()
	strategy := basic.New(&user.MockService{Err: errors.New("test")})

	body := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, userEmail, userPassword))
	req := httptest.NewRequest("GET", "/", body)
	_, err := strategy.Validate(req)
	assert.EqualError(t, err, "test")
	assert.False(t, strategy.IsErrCredentials(err))
	assert.False(t, strategy.IsErrDecoding(err))
}

func TestValidateFailsUserNotFound(t *testing.T) {
	t.Parallel()
	strategy := basic.New(&user.MockService{Err: user.ErrNotFound})

	body := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, userEmail, userPassword))
	req := httptest.NewRequest("GET", "/", body)
	_, err := strategy.Validate(req)
	assert.EqualError(t, err, "invalid email or password")
	assert.True(t, strategy.IsErrCredentials(err))
	assert.False(t, strategy.IsErrDecoding(err))
}
