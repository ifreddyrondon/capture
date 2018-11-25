package basic_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/pkg/auth/authentication/strategy/basic"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationPayloadValidUnmarshalJSON(t *testing.T) {
	t.Parallel()

	expected := basic.Credentials{
		Email:    "ifreddyrondon@example.com",
		Password: "b4KeHAYy3u9v=ZQX",
	}

	var model basic.Credentials
	err := model.UnmarshalJSON([]byte(`{"email":"ifreddyrondon@example.com", "password": "b4KeHAYy3u9v=ZQX"}`))
	require.Nil(t, err)
	assert.Equal(t, expected.Email, model.Email)
	assert.Equal(t, expected.Password, model.Password)
}

func TestAuthenticationPayloadInvalidUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		payload []byte
		errs    []string
	}{
		{
			"invalid payload",
			[]byte(`{`),
			[]string{"cannot unmarshal json into valid credentials"},
		},
		{
			"empty email and empty password",
			[]byte(`{"email":""}`),
			[]string{"email must not be blank", "password must not be blank"},
		},
		{
			"invalid email and empty password",
			[]byte(`{"email":"abc@abc."}`),
			[]string{"invalid email", "password must not be blank"},
		},
		{
			"empty password",
			[]byte(`{"email":"abc@example.com"}`),
			[]string{"password must not be blank"},
		},
		{
			"empty email",
			[]byte(`{"password":"b4KeHAYy3u9v=ZQX"}`),
			[]string{"email must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var model basic.Credentials
			err := model.UnmarshalJSON(tc.payload)
			assert.Error(t, err)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}
