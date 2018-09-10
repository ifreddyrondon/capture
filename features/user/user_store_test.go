package user_test

import (
	"sync"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/ifreddyrondon/capture/features/user"
)

var once sync.Once
var db *gorm.DB

func getDB(t *testing.T) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open("postgres", "postgres://localhost/captures_app_test?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}

func setupStore(t *testing.T) (user.Store, func()) {
	store := user.NewPGStore(getDB(t).Table("users"))
	store.Migrate()
	teardown := func() { store.Drop() }

	return store, teardown
}
