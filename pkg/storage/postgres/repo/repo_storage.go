package repo

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

type repoNotFound string

func (u repoNotFound) Error() string  { return string(u) }
func (u repoNotFound) NotFound() bool { return true }

// PGStorage postgres storage layer
type PGStorage struct{ db *pg.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *pg.DB) *PGStorage { return &PGStorage{db: db} }

// CreateSchema runs schema migration.
func (p *PGStorage) CreateSchema() error {
	opts := &orm.CreateTableOptions{IfNotExists: true}
	err := p.db.CreateTable(&domain.Repository{}, opts)
	if err != nil {
		return errors.Wrap(err, "creating repository schema")
	}
	return nil
}

// Drop delete schema.
func (p *PGStorage) Drop() error {
	opts := &orm.DropTableOptions{IfExists: true}
	err := p.db.DropTable(&domain.Repository{}, opts)
	if err != nil {
		return errors.Wrap(err, "dropping repository schema")
	}
	return nil
}

// Save capture into the database.
func (p *PGStorage) SaveRepo(repo *domain.Repository) error {
	if err := p.db.Insert(repo); err != nil {
		return errors.Wrap(err, "err saving repo with pgstorage")
	}
	return nil
}

func (p *PGStorage) List(l *domain.Listing) ([]domain.Repository, int64, error) {
	var repos []domain.Repository
	f := filter(*l)
	total, err := p.db.Model(&repos).Apply(f.Filter).SelectAndCount()
	if err != nil {
		return nil, 0, errors.Wrap(err, "err listing repo with pgstorage")
	}
	return repos, int64(total), nil
}

func (p *PGStorage) Get(id kallax.ULID) (*domain.Repository, error) {
	var repo domain.Repository
	err := p.db.Model(&repo).Where("id = ?", id).First()
	if err != nil {
		return nil, errors.WithStack(repoNotFound(fmt.Sprintf("repo with id %s not found", id)))
	}
	return &repo, nil
}
