package repo

import (
	"time"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/domain"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

type Repository struct {
	ID            kallax.ULID `sql:"type:uuid" gorm:"primary_key"`
	Name          string
	CurrentBranch string
	Visibility    string
	CreatedAt     time.Time `sql:"not null"`
	UpdatedAt     time.Time `sql:"not null"`
	DeletedAt     *time.Time
	UserID        kallax.ULID
}

// PGStorage postgres storage layer
type PGStorage struct{ db *gorm.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *gorm.DB) *PGStorage { return &PGStorage{db: db} }

// Migrate (panic) runs schema migration.
func (p *PGStorage) Migrate() {
	p.db.AutoMigrate(Repository{})
}

// Drop (panic) delete schema.
func (p *PGStorage) Drop() {
	p.db.DropTableIfExists(Repository{})
}

func getRepo(domainRepo *pkg.Repository) *Repository {
	return &Repository{
		ID:            domainRepo.ID,
		Name:          domainRepo.Name,
		CurrentBranch: domainRepo.CurrentBranch,
		Visibility:    string(domainRepo.Visibility),
		CreatedAt:     domainRepo.CreatedAt,
		UpdatedAt:     domainRepo.UpdatedAt,
		DeletedAt:     domainRepo.DeletedAt,
		UserID:        domainRepo.UserID,
	}
}

// Save capture into the database.
func (p *PGStorage) SaveRepo(domainRepo *pkg.Repository) error {
	r := getRepo(domainRepo)
	if err := p.db.Create(r).Error; err != nil {
		return errors.Wrap(err, "err saving repo with pgstorage")
	}
	return nil
}

func (p *PGStorage) List(l *domain.Listing) ([]pkg.Repository, error) {
	var results []pkg.Repository
	f := &pkg.Repository{}
	if l.Owner != "" {
		id, _ := kallax.NewULIDFromText(l.Owner)
		f.UserID = id
	}
	if l.Visibility != nil {
		f.Visibility = *l.Visibility
	}
	err := p.db.
		Where(f).
		Order(l.SortKey).
		Offset(l.Offset).
		Limit(l.Limit).
		Find(&results).Error
	if err != nil {
		return nil, errors.Wrap(err, "err listing repo with pgstorage")
	}
	return results, nil
}
