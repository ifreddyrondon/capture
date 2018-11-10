package user_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/user"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func setupService(t *testing.T) (*user.StoreService, func()) {
	toml := []byte(`PG="postgres://localhost/captures_app_test?sslmode=disable"`)
	cfg, err := config.New(config.Source(bytes.NewBuffer(toml)))
	if err != nil {
		t.Fatal(err)
	}

	db := cfg.Resources.Get("database").(*gorm.DB)
	store := user.NewPGStore(db.Table("users"))
	store.Migrate()

	return user.NewService(store), func() { store.Drop() }
}

func TestSaveUser(t *testing.T) {
	t.Parallel()

	store := &user.MockStore{}
	service := user.NewService(store)

	u := features.User{Email: "test@example.com"}
	err := service.Save(&u)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.NotNil(t, u.CreatedAt)
	assert.NotNil(t, u.UpdatedAt)
}

func TestErrWhenSaveUserWithSameEmail(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	u := features.User{Email: "test@example.com", ID: kallax.NewULID()}
	service.Save(&u)

	u2 := features.User{Email: "test@example.com", ID: kallax.NewULID()}
	err := service.Save(&u2)
	assert.Error(t, err)
	assert.Equal(t, "email 'test@example.com' already exists", err.Error())
}

func TestErrSaveUser(t *testing.T) {
	t.Parallel()

	store := &user.MockStore{Err: errors.New("test")}
	service := user.NewService(store)
	u := features.User{Email: "test@example.com"}
	err := service.Save(&u)
	assert.EqualError(t, err, "test")
}

func TestGetByEmail(t *testing.T) {
	t.Parallel()

	u := features.User{Email: "test@example.com"}
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

	u := features.User{Email: "test@example.com", ID: kallax.NewULID()}
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
