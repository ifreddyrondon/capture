package user

import (
	"fmt"

	"github.com/ifreddyrondon/capture/pkg"
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
	p.db.AutoMigrate(User{})
}

// Drop (panic) delete schema.
func (p *PGStorage) Drop() {
	p.db.DropTableIfExists(User{})
}

func getUser(domain *pkg.User) *User {
	return &User{
		ID:           domain.ID,
		Email:        domain.Email,
		Password:     domain.Password,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
		DeletedAt:    domain.DeletedAt,
		Repositories: domain.Repositories,
	}
}

func getDomainUser(u User) *pkg.User {
	return &pkg.User{
		ID:           u.ID,
		Email:        u.Email,
		Password:     u.Password,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
		Repositories: u.Repositories,
	}
}

// Save capture into the database.
func (p *PGStorage) SaveUser(user *pkg.User) error {
	u := getUser(user)
	err := p.db.Create(u).Error
	if err != nil {
		if isUniqueConstraintError(err, "uix_users_email") {
			return errors.WithStack(uniqueConstraintErr(err.Error()))
		}
		return errors.WithStack(err)
	}
	return nil
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStorage) GetUserByEmail(email string) (*pkg.User, error) {
	f := &User{Email: email}
	var result User
	if p.db.Where(f).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with email %s not found", email)))
	}
	return getDomainUser(result), nil
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStorage) GetUserByID(id kallax.ULID) (*pkg.User, error) {
	f := &User{ID: id}
	var result User
	if p.db.Where(f).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with id %v not found", id)))
	}
	return getDomainUser(result), nil
}
