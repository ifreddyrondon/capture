package user_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/database"
	"github.com/jinzhu/gorm"

	"github.com/ifreddyrondon/capture/user"
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

func setupStore(t *testing.T) (user.Store, func()) {
	store := user.NewPGStore(getDB().Table("users"))
	store.Migrate()
	teardown := func() { store.Drop() }

	return store, teardown
}
