package user_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/ifreddyrondon/capture/features/user"
	"github.com/ifreddyrondon/capture/internal/config"
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

func setupStore(t *testing.T) (user.Store, func()) {
	store := user.NewPGStore(getDB().Table("users"))
	store.Migrate()
	teardown := func() { store.Drop() }

	return store, teardown
}
