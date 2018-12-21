package user

import (
	"fmt"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type userNotFound string

func (u userNotFound) Error() string  { return string(u) }
func (u userNotFound) NotFound() bool { return true }

type uniqueConstraintErr string

func (u uniqueConstraintErr) Error() string          { return string(u) }
func (u uniqueConstraintErr) UniqueConstraint() bool { return true }

func isUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		fmt.Println(pqErr.Constraint)
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}

// PGStorage postgres storage layer
type PGStorage struct{ db *gorm.DB }

// NewPGStorage creates a new instance of PGStorage
func NewPGStorage(db *gorm.DB) *PGStorage { return &PGStorage{db: db} }

// Migrate (panic) runs schema migration.
func (p *PGStorage) Migrate() {
	p.db.AutoMigrate(domain.User{})
}

// Drop (panic) delete schema.
func (p *PGStorage) Drop() {
	p.db.DropTableIfExists(domain.User{})
}

// Save capture into the database.
func (p *PGStorage) SaveUser(user *domain.User) error {
	err := p.db.Create(user).Error
	if err != nil {
		if isUniqueConstraintError(err, "uix_users_email") {
			return errors.WithStack(uniqueConstraintErr(err.Error()))
		}
		return errors.WithStack(err)
	}
	return nil
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStorage) GetUserByEmail(email string) (*domain.User, error) {
	var result domain.User
	if p.db.Where(&domain.User{Email: email}).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with email %s not found", email)))
	}
	return &result, nil
}

// GetByEmail a user by id, if not found returns an error
func (p *PGStorage) GetUserByID(id kallax.ULID) (*domain.User, error) {
	var result domain.User
	if p.db.Where(&domain.User{ID: id}).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with id %v not found", id)))
	}
	return &result, nil
}
