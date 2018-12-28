package capture

import (
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// PGStorage postgres storage layer
type PGStorage struct{ db *gorm.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *gorm.DB) *PGStorage { return &PGStorage{db: db} }

// Migrate (panic) runs schema migration.
func (p *PGStorage) Migrate() {
	p.db.AutoMigrate(Capture{})
}

// Drop (panic) delete schema.
func (p *PGStorage) Drop() {
	p.db.DropTableIfExists(Capture{})
}

func (p *PGStorage) CreateCapture(capt *domain.Capture) error {
	c := getCapture(*capt)
	if err := p.db.Create(c).Error; err != nil {
		return errors.Wrap(err, "err saving capture with pgstorage")
	}
	return nil
}

func getCapture(c domain.Capture) *Capture {
	result := &Capture{
		ID:        c.ID,
		Payload:   payload(c.Payload),
		Tags:      pq.StringArray(c.Tags),
		Timestamp: c.Timestamp,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
	if c.Location != nil {
		result.Location = &point{
			LAT:       c.Location.LAT,
			LNG:       c.Location.LNG,
			Elevation: c.Location.Elevation,
		}
	}
	return result
}
