package repo

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/domain"

	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

type repoNotFound string

func (u repoNotFound) Error() string  { return string(u) }
func (u repoNotFound) NotFound() bool { return true }

// PGStorage postgres storage layer
type PGStorage struct{ db *gorm.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *gorm.DB) *PGStorage { return &PGStorage{db: db} }

// Migrate (panic) runs schema migration.
func (p *PGStorage) Migrate() {
	p.db.AutoMigrate(domain.Repository{})
}

// Drop (panic) delete schema.
func (p *PGStorage) Drop() {
	p.db.DropTableIfExists(domain.Repository{})
}

// Save capture into the database.
func (p *PGStorage) SaveRepo(repo *domain.Repository) error {
	if err := p.db.Create(repo).Error; err != nil {
		return errors.Wrap(err, "err saving repo with pgstorage")
	}
	return nil
}

func (p *PGStorage) List(l *domain.Listing) ([]domain.Repository, error) {
	var results []domain.Repository
	f := &domain.Repository{}
	if l.Owner != nil {
		f.UserID = *l.Owner
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

func (p *PGStorage) Get(id kallax.ULID) (*domain.Repository, error) {
	var result domain.Repository
	if p.db.Where(&domain.Repository{ID: id}).First(&result).RecordNotFound() {
		return nil, errors.WithStack(repoNotFound(fmt.Sprintf("repo with id %s not found", id)))
	}
	return &result, nil
}
