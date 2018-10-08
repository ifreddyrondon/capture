package user

import (
	"github.com/ifreddyrondon/capture/features"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save user into the database.
	Save(user *features.User) error
	// Get a user by email from database
	GetByEmail(string) (*features.User, error)
	// Get a user by id from database
	GetByID(kallax.ULID) (*features.User, error)
}

// PGStore implementation of user.Store for Postgres database.
type PGStore struct {
	db *gorm.DB
}

// NewPGStore creates a PGStore
func NewPGStore(db *gorm.DB) *PGStore {
	return &PGStore{db: db}
}

// Migrate (panic) runs schema migration.
func (p *PGStore) Migrate() {
	p.db.AutoMigrate(features.User{})
}

// Drop (panic) delete schema.
func (p *PGStore) Drop() {
	p.db.DropTableIfExists(features.User{})
}

// Save capture into the database.
func (p *PGStore) Save(user *features.User) error {
	return p.db.Create(user).Error
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStore) GetByEmail(email string) (*features.User, error) {
	var result features.User
	if p.db.Where(&features.User{Email: email}).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}

// GetByID a user by ID, if not found returns an error
func (p *PGStore) GetByID(id kallax.ULID) (*features.User, error) {
	var result features.User
	if p.db.Where(&features.User{ID: id}).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}

type MockStore struct {
	User *features.User
	Err  error
}

func (m *MockStore) Save(user *features.User) error                  { return m.Err }
func (m *MockStore) GetByEmail(email string) (*features.User, error) { return m.User, m.Err }
func (m *MockStore) GetByID(id kallax.ULID) (*features.User, error)  { return m.User, m.Err }
