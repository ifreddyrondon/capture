package user_test

import (
	"errors"
	"testing"

	"github.com/ifreddyrondon/capture/features/user"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func setupService(t *testing.T) (*user.StoreService, func()) {
	store, teardown := setupStore(t)
	return user.NewService(store), teardown
}

func TestSaveUser(t *testing.T) {
	t.Parallel()

	store := &user.MockStore{}
	service := user.NewService(store)

	u := user.User{Email: "test@example.com"}
	err := service.Save(&u)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.NotNil(t, u.CreatedAt)
	assert.NotNil(t, u.UpdatedAt)
}

func TestErrWhenSaveUserWithSameEmail(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	u := user.User{Email: "test@example.com"}
	service.Save(&u)

	u2 := user.User{Email: "test@example.com"}
	err := service.Save(&u2)
	assert.Error(t, err)
	assert.Equal(t, "email 'test@example.com' already exists", err.Error())
}

func TestErrSaveUser(t *testing.T) {
	t.Parallel()

	store := &user.MockStore{Err: errors.New("test")}
	service := user.NewService(store)
	u := user.User{Email: "test@example.com"}
	err := service.Save(&u)
	assert.EqualError(t, err, "test")
}

func TestGetByEmail(t *testing.T) {
	t.Parallel()

	u := user.User{Email: "test@example.com"}
	store := &user.MockStore{User: &u}
	service := user.NewService(store)

	tempUser, err := service.GetByEmail("test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, u.ID, tempUser.ID)
}

func TestGetByEmailMissing(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	_, err := service.GetByEmail("test@example.com")
	assert.Error(t, err)
	assert.Equal(t, user.ErrNotFound, err)
}

func TestGetByID(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	u := user.User{Email: "test@example.com"}
	service.Save(&u)

	tempUser, err := service.GetByID(u.ID)
	assert.Nil(t, err)
	assert.Equal(t, u.Email, tempUser.Email)
}

func TestGetByIDMissing(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	_, err := service.GetByID(kallax.NewULID())
	assert.Error(t, err)
	assert.Equal(t, user.ErrNotFound, err)
}
