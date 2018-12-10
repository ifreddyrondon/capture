package user

import (
	"fmt"
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"
)

type User struct {
	ID           kallax.ULID `sql:"type:uuid" gorm:"primary_key"`
	Email        string      `sql:"not null" gorm:"unique_index"`
	Password     []byte
	CreatedAt    time.Time `sql:"not null"`
	UpdatedAt    time.Time `sql:"not null"`
	DeletedAt    *time.Time
	Repositories []pkg.Repository `gorm:"ForeignKey:UserID"`
}

type userNotFound string

func (u userNotFound) Error() string  { return string(u) }
func (u userNotFound) NotFound() bool { return true }

type uniqueConstraintErr string

func (u uniqueConstraintErr) Error() string          { return string(u) }
func (u uniqueConstraintErr) UniqueConstraint() bool { return true }

type invalidIDErr string

func (i invalidIDErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidIDErr) IsInvalid() bool { return true }

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

func getUserId(idStr string) (kallax.ULID, error) {
	id, err := kallax.NewULIDFromText(idStr)
	if err != nil {
		return kallax.ULID{}, invalidIDErr(fmt.Sprintf("%v is not a valid ULID", idStr))
	}
	return id, nil
}

func getUser(domain *pkg.User) (*User, error) {
	id, err := getUserId(domain.ID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ulid format when getting user")
	}
	return &User{
		ID:           id,
		Email:        domain.Email,
		Password:     domain.Password,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
		DeletedAt:    domain.DeletedAt,
		Repositories: domain.Repositories,
	}, nil
}

func getDomainUser(u User) *pkg.User {
	return &pkg.User{
		ID:           u.ID.String(),
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
	u, err := getUser(user)
	if err != nil {
		return err
	}
	err = p.db.Create(u).Error
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
func (p *PGStorage) GetUserByID(idStr string) (*pkg.User, error) {
	id, err := getUserId(idStr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ulid format when GetUserByID")
	}
	f := &User{ID: id}
	var result User
	if p.db.Where(f).First(&result).RecordNotFound() {
		return nil, errors.WithStack(userNotFound(fmt.Sprintf("user with id %v not found", id)))
	}
	return getDomainUser(result), nil
}
