package user_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/gocapture/database"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/gocapture/user"
)

var once sync.Once
var db *gorm.DB

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setupService(t *testing.T) (*user.PGService, func()) {
	service := user.PGService{DB: getDB().Table("users")}
	service.Migrate()
	teardown := func() { service.Drop() }

	return &service, teardown
}

func TestSaveUser(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

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
