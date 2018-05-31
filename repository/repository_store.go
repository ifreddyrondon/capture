package repository

import "github.com/jinzhu/gorm"

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save a repository.
	Save(*Repository) error
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
	pgs.db.AutoMigrate(Repository{})
}

// Drop (panic) delete the repository schema.
func (pgs *PGStore) Drop() {
	pgs.db.DropTableIfExists(Repository{})
}

// Save a repository into the database.
func (pgs *PGStore) Save(r *Repository) error {
	return pgs.db.Create(r).Error
}
