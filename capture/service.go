package capture

import "github.com/jinzhu/gorm"

// Service is the interface implemented by capture
// It make CRUD operations over captures.
type Service interface {
	// Save capture into the database.
	Save(*Capture) (*Capture, error)
	// SaveBulk captures into the database.
	SaveBulk(...*Capture) (Captures, error)
	// List retrieve the count captures from start index.
	List(start, count int) (Captures, error)
	// Get a capture by id
	Get(uint64) (*Capture, error)
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
func (pgs *PGService) Save(capt *Capture) (*Capture, error) {
	if err := pgs.DB.Create(capt).Error; err != nil {
		return nil, err
	}
	return capt, nil
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
func (pgs *PGService) Get(id uint64) (*Capture, error) {
	var result Capture
	if pgs.DB.First(&result, id).RecordNotFound() {
		return nil, ErrorNotFound
	}
	return &result, nil
}

// Delete a capture by id
func (pgs *PGService) Delete(capt *Capture) error {
	if err := pgs.DB.Delete(capt).Error; err != nil {
		return err
	}
	return nil
}

// Update a capture
func (pgs *PGService) Update(original *Capture, updates Capture) error {
	if err := pgs.DB.Model(original).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
