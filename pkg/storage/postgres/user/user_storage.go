package user

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type userNotFound string

func (u userNotFound) Error() string  { return string(u) }
func (u userNotFound) NotFound() bool { return true }

type uniqueConstraintErr string

func (u uniqueConstraintErr) Error() string          { return string(u) }
func (u uniqueConstraintErr) UniqueConstraint() bool { return true }

func isUniqueConstraintError(err error) bool {
	if pqErr, ok := err.(pg.Error); ok {
		return pqErr.IntegrityViolation()
	}
	return false
}

// PGStorage postgres storage layer
type PGStorage struct{ db *pg.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *pg.DB) *PGStorage { return &PGStorage{db: db} }

// CreateSchema runs schema migration.
func (p *PGStorage) CreateSchema() error {
	opts := &orm.CreateTableOptions{IfNotExists: true}
	err := p.db.CreateTable(&domain.User{}, opts)
	if err != nil {
		return errors.Wrap(err, "creating user schema")
	}
	return nil
}

// Drop delete schema.
func (p *PGStorage) Drop() error {
	opts := &orm.DropTableOptions{IfExists: true}
	err := p.db.DropTable(&domain.User{}, opts)
	if err != nil {
		return errors.Wrap(err, "dropping user schema")
	}
	return nil
}

// Save capture into the database.
func (p *PGStorage) SaveUser(user *domain.User) error {
	err := p.db.Insert(user)
	if err != nil {
		if isUniqueConstraintError(err) {
			return errors.WithStack(uniqueConstraintErr(err.Error()))
		}
		return errors.WithStack(err)
	}
	return nil
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStorage) GetUserByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := p.db.Model(&u).Where("email = ?", email).First()
	if err != nil {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with email %s not found", email)))
	}
	return &u, nil
}

// GetByEmail a user by id, if not found returns an error
func (p *PGStorage) GetUserByID(id kallax.ULID) (*domain.User, error) {
	var u domain.User
	err := p.db.Model(&u).Where("id = ?", id).First()
	if err != nil {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with id %v not found", id)))
	}
	return &u, nil
}
