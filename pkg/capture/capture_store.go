package capture

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save capture into the database.
	Save(*pkg.Capture) error
	// SaveBulk captures into the database.
	SaveBulk(...*pkg.Capture) (Captures, error)
	// List retrieve captures from start index to count.
	List(start, count int) (Captures, error)
	// Get a capture by id
	Get(kallax.ULID) (*pkg.Capture, error)
	// Delete a capture by id
	Delete(*pkg.Capture) error
	// Update a capture from an updated one, will only update those changed & non blank fields.
	Update(original *pkg.Capture, updates pkg.Capture) error
}

// PGStore implementation of capture.Store for Postgres database.
type PGStore struct {
	db *gorm.DB
}

// NewPGStore creates a PGStore
func NewPGStore(db *gorm.DB) *PGStore {
	return &PGStore{db: db}
}

// Migrate (panic) runs schema migration.
func (pgs *PGStore) Migrate() {
	pgs.db.AutoMigrate(pkg.Capture{})
}

// Drop (panic) delete schema.
func (pgs *PGStore) Drop() {
	pgs.db.DropTableIfExists(pkg.Capture{})
}

// Save capture into the database.
func (pgs *PGStore) Save(capt *pkg.Capture) error {
	return pgs.db.Create(capt).Error
}

// SaveBulk captures into the database.
func (pgs *PGStore) SaveBulk(captures ...*pkg.Capture) (Captures, error) {
	// TODO: bash create
	for _, c := range captures {
		if err := pgs.db.Create(c).Error; err != nil {
			continue
		}
	}
	return captures, nil
}

// List retrieve the count captures from start index.
func (pgs *PGStore) List(start, count int) (Captures, error) {
	results := Captures{}
	if err := pgs.db.Order("updated_at").Offset(start).Limit(count).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Get a capture by id
func (pgs *PGStore) Get(id kallax.ULID) (*pkg.Capture, error) {
	var result pkg.Capture
	if pgs.db.Where(&pkg.Capture{ID: id}).First(&result).RecordNotFound() {
		return nil, ErrorNotFound
	}
	return &result, nil
}

// Delete a capture by id
func (pgs *PGStore) Delete(capt *pkg.Capture) error {
	return pgs.db.Delete(capt).Error
}

// Update a capture
func (pgs *PGStore) Update(original *pkg.Capture, updates pkg.Capture) error {
	return pgs.db.Model(original).Updates(updates).Error
}
