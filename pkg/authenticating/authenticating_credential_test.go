package authenticating_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/bastion/binder"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/authenticating"
)

func TestValidateUserOK(t *testing.T) {
	t.Parallel()

	email := "ifreddyrondon@example.com"
	pass := "b4KeHAYy3u9v=ZQX"

	body := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, email, pass))
	expected := authenticating.BasicCredential{Email: email, Password: pass}

	r, _ := http.NewRequest("POST", "/", body)

	var credentials authenticating.BasicCredential
	err := binder.JSON.FromReq(r, &credentials)
	assert.Nil(t, err)
	assert.Equal(t, expected.Email, credentials.Email)
	assert.Equal(t, expected.Password, credentials.Password)
}

func TestValidateUserError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		body string
		errs []string
	}{
		{
			"invalid payload",
			`{`,
			[]string{"cannot unmarshal json body"},
		},
		{
			"empty email and empty password",
			`{"email":""}`,
			[]string{"email must not be blank", "password must not be blank"},
		},
		{
			"invalid email and empty password",
			`{"email":"abc@abc."}`,
			[]string{"invalid email", "password must not be blank"},
		},
		{
			"empty password",
			`{"email":"abc@example.com"}`,
			[]string{"password must not be blank"},
		},
		{
			"empty email",
			`{"password":"b4KeHAYy3u9v=ZQX"}`,
			[]string{"email must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))
			var credentials authenticating.BasicCredential
			err := binder.JSON.FromReq(r, &credentials)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}
