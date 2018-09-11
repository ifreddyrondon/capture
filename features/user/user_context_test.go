package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/features/user"
)

func TestContextManagerGetUserOK(t *testing.T) {
	ctx := context.Background()
	u := user.User{ID: kallax.NewULID(), Email: "test@example.com"}
	ctx = user.WithUser(ctx, &u)

	u2, err := user.GetUser(ctx)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, u2.ID)
	assert.Equal(t, u.Email, u2.Email)
}

func TestContextManagerGetUserMissingUser(t *testing.T) {
	ctx := context.Background()
	_, err := user.GetUser(ctx)
	assert.EqualError(t, err, "user not found in context")
}

func TestContextManagerGetUserIDOK(t *testing.T) {
	ctx := context.Background()

	u := user.User{ID: kallax.NewULID(), Email: "test@example.com"}
	ctx = user.WithUser(ctx, &u)

	id, err := user.GetUserID(ctx)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, id)
}

func TestContextManagerGetUserIDMissingUser(t *testing.T) {
	ctx := context.Background()

	_, err := user.GetUserID(ctx)
	assert.EqualError(t, err, "user not found in context")
}
