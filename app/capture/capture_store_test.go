package capture_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/app/capture"
	"github.com/ifreddyrondon/capture/database"
	"github.com/jinzhu/gorm"
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

func setupStore(t *testing.T) (capture.Store, func()) {
	store := capture.NewPGStore(getDB().Table("captures"))
	store.Migrate()
	teardown := func() { store.Drop() }

	return store, teardown
}
