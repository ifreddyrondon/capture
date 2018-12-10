package repo

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
)

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
	}
}

// Save capture into the database.
func (p *PGStorage) SaveUser(domainRepo *pkg.Repository) error {
	r := getRepo(domainRepo)
	return p.db.Create(r).Error
}
