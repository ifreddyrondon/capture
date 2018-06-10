package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/app/user"
)

func TestContextManagerGetUserOK(t *testing.T) {
	ctxManager := user.NewContextManager()
	ctx := context.Background()

	u := user.User{ID: kallax.NewULID(), Email: "test@example.com"}
	ctx = ctxManager.WithUser(ctx, &u)

	u2, err := ctxManager.GetUser(ctx)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, u2.ID)
	assert.Equal(t, u.Email, u2.Email)
}

func TestContextManagerGetUserMissingUser(t *testing.T) {
	ctxManager := user.NewContextManager()
	ctx := context.Background()

	_, err := ctxManager.GetUser(ctx)
	assert.EqualError(t, err, "user not found in context")
}

func TestContextManagerGetUserIDOK(t *testing.T) {
	ctxManager := user.NewContextManager()
	ctx := context.Background()

	u := user.User{ID: kallax.NewULID(), Email: "test@example.com"}
	ctx = ctxManager.WithUser(ctx, &u)

	id, err := ctxManager.GetUserID(ctx)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, id)
}

func TestContextManagerGetUserIDMissingUser(t *testing.T) {
	ctxManager := user.NewContextManager()
	ctx := context.Background()

	_, err := ctxManager.GetUserID(ctx)
	assert.EqualError(t, err, "user not found in context")
}
