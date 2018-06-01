package repository_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/app/repository"
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

func setupStore(t *testing.T) (repository.Store, func()) {
	repo := repository.NewPGStore(getDB().Table("repositories"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	return repo, teardown
}
