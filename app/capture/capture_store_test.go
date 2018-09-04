package capture_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/app/capture"
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

func setupStore(t *testing.T) (capture.Store, func()) {
	store := capture.NewPGStore(getDB().Table("captures"))
	store.Migrate()
	teardown := func() { store.Drop() }

	return store, teardown
}
