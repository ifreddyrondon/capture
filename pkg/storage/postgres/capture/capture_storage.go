package capture

import (
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type captureNotFound string

func (u captureNotFound) Error() string  { return string(u) }
func (u captureNotFound) NotFound() bool { return true }

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

func (p *PGStorage) List(l *domain.Listing) ([]domain.Capture, error) {
	var captures []domain.Capture
	f := filter(*l)
	err := p.db.Model(&captures).Apply(f.Filter).Select()
	if err != nil {
		return nil, errors.Wrap(err, "err listing captures with pgstorage")
	}
	return captures, nil
}

func (p *PGStorage) Get(captureID, repoID kallax.ULID) (*domain.Capture, error) {
	var capt domain.Capture
	err := p.db.Model(&capt).
		Where("id = ?", captureID).
		Where("repository_id = ?", repoID).
		First()
	if err != nil {
		errStr := fmt.Sprintf("capture with id %s not found in repo %v", captureID, repoID)
		return nil, errors.WithStack(captureNotFound(errStr))
	}
	return &capt, nil
}

func (p *PGStorage) Delete(capt *domain.Capture) error {
	t := time.Now()
	capt.DeletedAt = &t
	if err := p.db.Update(capt); err != nil {
		errStr := fmt.Sprintf("error deleting capture %s in repo %v", capt.ID, capt.RepositoryID)
		return errors.Wrap(err, errStr)
	}
	return nil
}
