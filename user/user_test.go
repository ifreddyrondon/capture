package user_test

import (
	"encoding/json"
	"testing"
	"time"

	kallax "gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/gocapture/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalValidUser(t *testing.T) {
	t.Parallel()

	expected := user.User{
		Email: "ifreddyrondon@gmail.com",
	}

	result := user.User{}
	err := result.UnmarshalJSON([]byte(`{"email":"ifreddyrondon@gmail.com"}`))
	require.Nil(t, err)
	assert.Equal(t, expected.Email, result.Email)
}

func TestUnmarshalInValidUser(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		payload []byte
		err     string
	}{
		{
			"invalid payload",
			[]byte(`{`),
			"cannot unmarshal json into valid user",
		},
		{
			"empty email",
			[]byte(`{"email":""}`),
			"email required",
		},
		{
			"invalid email - abc@abc.",
			[]byte(`{"email":"abc@abc."}`),
			"invalid email",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := user.User{}
			err := result.UnmarshalJSON(tc.payload)
			require.NotNil(t, err)
			assert.Error(t, err, tc.err)
		})
	}
}

func TestMarshalUser(t *testing.T) {
	t.Parallel()
	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := []byte(`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","email":"test@test.com","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`)
	user := user.User{
		Email: "test@test.com",
		ID: func() kallax.ULID {
			id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			return id
		}(),
		CreatedAt: d,
		UpdatedAt: d,
	}

	result, err := json.Marshal(user)
	require.Nil(t, err)
	assert.Equal(t, expected, result)
}
