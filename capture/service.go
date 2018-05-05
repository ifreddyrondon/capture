package capture

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

// Service is the interface implemented by capture
// It make CRUD operations over captures.
type Service interface {
	// Save capture into the database.
	Save(*Capture) error
	// SaveBulk captures into the database.
	SaveBulk(...*Capture) (Captures, error)
	// List retrieve captures from start index to count.
	List(start, count int) (Captures, error)
	// Get a capture by id
	Get(kallax.ULID) (*Capture, error)
	// Delete a capture by id
	Delete(*Capture) error
	// Update a capture from an updated one, will only update those changed & non blank fields.
	Update(original *Capture, updates Capture) error
}

// PGService implementation of capture.Service for Postgres database.
type PGService struct {
	DB *gorm.DB
}

// Migrate (panic) runs schema migration.
func (pgs *PGService) Migrate() {
	pgs.DB.AutoMigrate(Capture{})
}

// Drop (panic) delete schema.
func (pgs *PGService) Drop() {
	pgs.DB.DropTableIfExists(Capture{})
}

// Save capture into the database.
func (pgs *PGService) Save(capt *Capture) error {
	capt.ID = kallax.NewULID()
	return pgs.DB.Create(capt).Error
}

// SaveBulk captures into the database.
func (pgs *PGService) SaveBulk(captures ...*Capture) (Captures, error) {
	// TODO: bash create
	for _, c := range captures {
		if err := pgs.DB.Create(c).Error; err != nil {
			continue
		}
	}
	return captures, nil
}

// List retrieve the count captures from start index.
func (pgs *PGService) List(start, count int) (Captures, error) {
	results := Captures{}
	if err := pgs.DB.Order("updated_at").Offset(start).Limit(count).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Get a capture by id
func (pgs *PGService) Get(id kallax.ULID) (*Capture, error) {
	var result Capture
	if pgs.DB.Where(&Capture{ID: id}).First(&result).RecordNotFound() {
		return nil, ErrorNotFound
	}
	return &result, nil
}

// Delete a capture by id
func (pgs *PGService) Delete(capt *Capture) error {
	return pgs.DB.Delete(capt).Error
}

// Update a capture
func (pgs *PGService) Update(original *Capture, updates Capture) error {
	return pgs.DB.Model(original).Updates(updates).Error
}
