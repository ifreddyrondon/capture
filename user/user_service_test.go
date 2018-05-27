package user_test

import (
	"errors"
	"testing"

	"github.com/ifreddyrondon/gocapture/user"
	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*user.REPOService, func()) {
	repo, teardown := setupRepository(t)
	return user.NewService(repo), teardown
}

func TestSaveUser(t *testing.T) {
	t.Parallel()

	repo := user.NewFakeRepositoryDefaultPanic()
	repo.SaveHook = func(*user.User) error { return nil }
	service := user.NewService(repo)

	u := user.User{Email: "test@localhost.com"}
	err := service.Save(&u)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.NotNil(t, u.CreatedAt)
	assert.NotNil(t, u.UpdatedAt)
}

func TestErrWhenSaveUserWithSameEmail(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	u := user.User{Email: "test@localhost.com"}
	service.Save(&u)

	u2 := user.User{Email: "test@localhost.com"}
	err := service.Save(&u2)
	assert.Error(t, err)
	assert.Equal(t, "email 'test@localhost.com' already exists", err.Error())
}

func TestErrSaveUser(t *testing.T) {
	t.Parallel()

	repo := user.NewFakeRepositoryDefaultPanic()
	repo.SaveHook = func(*user.User) error { return errors.New("test") }
	service := user.NewService(repo)

	u := user.User{Email: "test@localhost.com"}
	err := service.Save(&u)
	assert.EqualError(t, err, "test")
}

func TestGetUser(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	u := user.User{Email: "test@localhost.com"}
	service.Save(&u)

	tempUser, err := service.Get("test@localhost.com")
	assert.Nil(t, err)
	assert.Equal(t, u.ID, tempUser.ID)
}

func TestGetMissingUser(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	_, err := service.Get("test@localhost.com")
	assert.Error(t, err)
	assert.Equal(t, user.ErrNotFound, err)
}
