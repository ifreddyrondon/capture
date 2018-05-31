package capture_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
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

func setupRepository(t *testing.T) (capture.Repository, func()) {
	repo := capture.NewPGRepository(getDB().Table("captures"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	return repo, teardown
}
