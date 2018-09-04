package repository_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/app/repository"
	"github.com/ifreddyrondon/capture/internal/config"
	"github.com/jinzhu/gorm"
)

var once sync.Once
var db *gorm.DB

func getDB() *gorm.DB {
	once.Do(func() {
		src := []byte(`PG="postgres://localhost/captures_app_test?sslmode=disable"`)
		cfg, _ := config.New(config.Source(bytes.NewBuffer(src)))
		db = cfg.Database
	})
	return db
}

func setupStore(t *testing.T) (repository.Store, func()) {
	repo := repository.NewPGStore(getDB().Table("repositories"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	return repo, teardown
}
