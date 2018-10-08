package features_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/user/decoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestUserPassword(t *testing.T) {
	t.Parallel()

	email, password := "test@localhost.com", "b4KeHAYy3u9v=ZQX"
	var u features.User
	err := decoder.User(decoder.PostUser{Email: &email, Password: &password}, &u)
	assert.Nil(t, err)
	assert.True(t, u.CheckPassword("b4KeHAYy3u9v=ZQX"))
	assert.False(t, u.CheckPassword("1"))
}

func TestMarshalUser(t *testing.T) {
	t.Parallel()
	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := []byte(`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","email":"test@example.com","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z"}`)
	u := features.User{
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
