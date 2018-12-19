package repo

import (
	"fmt"
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

type repoNotFound string

func (u repoNotFound) Error() string  { return string(u) }
func (u repoNotFound) NotFound() bool { return true }

type invalidIDErr string

func (i invalidIDErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidIDErr) IsInvalid() bool { return true }

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
		id, err := kallax.NewULIDFromText(l.Owner)
		if err != nil {
			return nil, invalidIDErr(fmt.Sprintf("%v is not a valid owner id", l.Owner))
		}
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

func (p *PGStorage) Get(idStr string) (*pkg.Repository, error) {
	var result pkg.Repository
	id, err := kallax.NewULIDFromText(idStr)
	if err != nil {
		return nil, invalidIDErr(fmt.Sprintf("%v is not a valid ULID", idStr))
	}
	if p.db.Where(&pkg.Repository{ID: id}).First(&result).RecordNotFound() {
		return nil, errors.WithStack(repoNotFound(fmt.Sprintf("repo with id %s not found", idStr)))
	}
	return &result, nil
}
