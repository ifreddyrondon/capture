package user_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/gocapture/database"
	"github.com/jinzhu/gorm"

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

func setupRepository(t *testing.T) (user.Repository, func()) {
	repo := user.NewPGRepository(getDB().Table("users"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	return repo, teardown
}
