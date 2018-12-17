package repository

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Get a repo by id
	Get(string) (*pkg.Repository, error)
	// Drop register if it is exist
	Drop()
}

// PGStore implementation of repository.Store for Postgres database.
type PGStore struct {
	db *gorm.DB
}

// NewPGStore creates a PGStore
func NewPGStore(db *gorm.DB) *PGStore {
	return &PGStore{db: db}
}

// Migrate (panic) runs schema migration for repository table.
func (pgs *PGStore) Migrate() {
	pgs.db.AutoMigrate(pkg.Repository{})
}

// Drop (panic) delete the repository schema.
func (pgs *PGStore) Drop() {
	pgs.db.DropTableIfExists(pkg.Repository{})
}

func (pgs *PGStore) Get(idStr string) (*pkg.Repository, error) {
	var result pkg.Repository
	// FIXME: handler id err
	id, _ := kallax.NewULIDFromText(idStr)
	if pgs.db.Where(&pkg.Repository{ID: id}).First(&result).RecordNotFound() {
		return nil, ErrorNotFound
	}
	return &result, nil
}
