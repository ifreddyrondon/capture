package user_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features/user/decoder"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/features/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPassword(t *testing.T) {
	t.Parallel()

	email, password := "test@localhost.com", "b4KeHAYy3u9v=ZQX"
	u, err := user.FromPostUser(decoder.PostUser{Email: &email, Password: &password})
	assert.Nil(t, err)
	assert.True(t, u.CheckPassword("b4KeHAYy3u9v=ZQX"))
	assert.False(t, u.CheckPassword("1"))
}

func TestMarshalUser(t *testing.T) {
	t.Parallel()
	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := []byte(`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","email":"test@example.com","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`)
	u := user.User{
		Email: "test@example.com",
		ID: func() kallax.ULID {
			id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			return id
		}(),
		CreatedAt: d,
		UpdatedAt: d,
	}

	result, err := json.Marshal(u)
	require.Nil(t, err)
	assert.Equal(t, expected, result)
}
