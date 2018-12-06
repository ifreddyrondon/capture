package postgres

import (
	"fmt"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type userNotFound string

func (u userNotFound) Error() string  { return string(u) }
func (u userNotFound) NotFound() bool { return true }

// PGStorage postgres storage layer
type PGStorage struct{ db *gorm.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *gorm.DB) *PGStorage { return &PGStorage{db: db} }

// Migrate (panic) runs schema migration.
func (p *PGStorage) Migrate() {
	p.db.AutoMigrate(pkg.User{})
}

// Drop (panic) delete schema.
func (p *PGStorage) Drop() {
	p.db.DropTableIfExists(pkg.User{})
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStorage) GetUserByEmail(email string) (*pkg.User, error) {
	f := &pkg.User{Email: email}
	var result pkg.User
	if p.db.Where(f).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with email %s not found", email)))
	}
	return &result, nil
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStorage) GetUserByID(id kallax.ULID) (*pkg.User, error) {
	f := &pkg.User{ID: id}
	var result pkg.User
	if p.db.Where(f).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with id %v not found", id)))
	}
	return &result, nil
}
