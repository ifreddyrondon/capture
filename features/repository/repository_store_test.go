package repository_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/features/repository"
	"github.com/jinzhu/gorm"
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

func setupStore(t *testing.T) (repository.Store, func()) {
	repo := repository.NewPGStore(getDB(t).Table("repositories"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	return repo, teardown
}
