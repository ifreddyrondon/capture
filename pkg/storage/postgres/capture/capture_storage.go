package capture

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

// PGStorage postgres storage layer
type PGStorage struct{ db *pg.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *pg.DB) *PGStorage { return &PGStorage{db: db} }

// CreateSchema runs schema migration.
func (p *PGStorage) CreateSchema() error {
	opts := &orm.CreateTableOptions{IfNotExists: true}
	err := p.db.CreateTable(&domain.Capture{}, opts)
	if err != nil {
		return errors.Wrap(err, "creating capture schema")
	}
	return nil
}

// Drop delete schema.
func (p *PGStorage) Drop() error {
	opts := &orm.DropTableOptions{IfExists: true}
	err := p.db.DropTable(&domain.Capture{}, opts)
	if err != nil {
		return errors.Wrap(err, "dropping capture schema")
	}
	return nil
}

func (p *PGStorage) CreateCapture(c *domain.Capture) error {
	if err := p.db.Insert(c); err != nil {
		return errors.Wrap(err, "err saving capture with pgstorage")
	}
	return nil
}
