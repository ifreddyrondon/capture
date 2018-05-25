package collection

import "github.com/jinzhu/gorm"

// Repository is the interface to be implemented by collection
// It make CRUD operations over a repository.
type Repository interface {
	// Save a collection.
	Save(*Collection) error
}

// PGRepository implementation of collection.Repository for Postgres database.
type PGRepository struct {
	db *gorm.DB
}

// NewPGRepository creates a PGRepository
func NewPGRepository(db *gorm.DB) *PGRepository {
	return &PGRepository{db: db}
}

// Migrate (panic) runs schema migration.
func (pgs *PGRepository) Migrate() {
	pgs.db.AutoMigrate(Collection{})
}

// Drop (panic) delete schema.
func (pgs *PGRepository) Drop() {
	pgs.db.DropTableIfExists(Collection{})
}

// Save capture into the database.
func (pgs *PGRepository) Save(c *Collection) error {
	return pgs.db.Create(c).Error
}
